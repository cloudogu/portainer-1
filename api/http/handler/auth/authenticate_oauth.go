package auth

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/cloudogu/portainer-ce/api"
	bolterrors "github.com/cloudogu/portainer-ce/api/bolt/errors"
	httperrors "github.com/cloudogu/portainer-ce/api/http/errors"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"log"
	"net/http"
)

type oauthPayload struct {
	Code string
}

func (payload *oauthPayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Code) {
		return errors.New("Invalid OAuth authorization code")
	}
	return nil
}

func (handler *Handler) authenticateOAuth(code string, settings *portainer.OAuthSettings) (portainer.OAuthUserData, error) {
	if code == "" {
		return portainer.OAuthUserData{}, errors.New("Invalid OAuth authorization code")
	}

	if settings == nil {
		return portainer.OAuthUserData{}, errors.New("Invalid OAuth configuration")
	}

	userData, err := handler.OAuthService.Authenticate(code, settings)
	if err != nil {
		return portainer.OAuthUserData{}, err
	}

	return userData, nil
}

func (handler *Handler) validateOAuth(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload oauthPayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	settings, err := handler.DataStore.Settings().Settings()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve settings from the database", err}
	}

	if settings.AuthenticationMethod != 3 {
		return &httperror.HandlerError{http.StatusForbidden, "OAuth authentication is not enabled", errors.New("OAuth authentication is not enabled")}
	}

	userData, err := handler.authenticateOAuth(payload.Code, &settings.OAuthSettings)
	if err != nil {
		log.Printf("[DEBUG] - OAuth authentication error: %s", err)
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to authenticate through OAuth", httperrors.ErrUnauthorized}
	}

	return handler.userProvisioning(w, &userData, settings)
}

func (handler *Handler) userProvisioning(w http.ResponseWriter, userData *portainer.OAuthUserData, settings *portainer.Settings) *httperror.HandlerError {
	user, err := handler.DataStore.User().UserByUsername(userData.Username)
	if err != nil && err != bolterrors.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve a user with the specified username from the database", err}
	}

	if user == nil && !settings.OAuthSettings.OAuthAutoCreateUsers {
		return &httperror.HandlerError{http.StatusForbidden, "Account not created beforehand in Portainer and automatic user provisioning not enabled", httperrors.ErrUnauthorized}
	}

	if user == nil {
		user, err = handler.createUser(userData)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist the new user inside the database", err}
		}
	}

	// Get new groups from profile
	newTeams, handlerErr := handler.getUserGroups(userData)
	if handlerErr != nil {
		return handlerErr
	}

	// Get old groups from database
	oldTeams, err := handler.DataStore.TeamMembership().TeamMembershipsByUserID(user.ID)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve team subscriptions from the database", err}
	}

	// Remove old groups that are not part of the new groups
	handlerErr = handler.removeUserFromOldGroups(&oldTeams, newTeams)
	if handlerErr != nil {
		return handlerErr
	}

	// Add new groups that are not part of the old groups
	handlerErr = handler.addUserToNewGroups(user, &oldTeams, newTeams)
	if handlerErr != nil {
		return handlerErr
	}

	// Handle admin group
	handlerErr = handler.checkAdminPrivileges(user, userData, settings)
	if handlerErr != nil {
		return handlerErr
	}

	user.OAuthToken = userData.OAuthToken
	return handler.writeToken(w, user)
}

func (handler *Handler) createUser(userData *portainer.OAuthUserData) (*portainer.User, error) {
	user := &portainer.User{
		Username: userData.Username,
		Role:     portainer.StandardUserRole,
	}

	err := handler.DataStore.User().CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (handler *Handler) getUserGroups(userData *portainer.OAuthUserData) (*[]portainer.Team, *httperror.HandlerError) {
	var newTeams []portainer.Team
	for _, teamName := range userData.Teams {
		team, err := handler.DataStore.Team().TeamByName(teamName)
		if err == bolterrors.ErrObjectNotFound {
			team = &portainer.Team{
				Name: teamName,
			}

			err = handler.DataStore.Team().CreateTeam(team)
			if err != nil {
				return nil, &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist the team inside the database", err}
			}
		} else if err != nil {
			return nil, &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve the team from the database", err}
		}
		newTeams = append(newTeams, *team)
	}
	return &newTeams, nil
}

func (handler *Handler) removeUserFromOldGroups(oldTeams *[]portainer.TeamMembership, newTeams *[]portainer.Team) *httperror.HandlerError {
	for _, oldTeam := range *oldTeams {
		var removeOldGroup = true
		for _, newTeam := range *newTeams {
			if oldTeam.TeamID == newTeam.ID {
				removeOldGroup = false
				break
			}
		}

		if removeOldGroup {
			err := handler.DataStore.TeamMembership().DeleteTeamMembership(oldTeam.ID)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "Unable to delete team subscriptions from the database", err}
			}
		}
	}
	return nil
}

func (handler *Handler) addUserToNewGroups(user *portainer.User, oldTeams *[]portainer.TeamMembership, newTeams *[]portainer.Team) *httperror.HandlerError {
	for _, newTeam := range *newTeams {
		var addToGroup = true
		for _, oldTeam := range *oldTeams {
			if oldTeam.TeamID == newTeam.ID {
				addToGroup = false
			}
		}

		if addToGroup {
			membership := &portainer.TeamMembership{
				UserID: user.ID,
				TeamID: newTeam.ID,
				Role:   portainer.TeamMember,
			}

			err := handler.DataStore.TeamMembership().CreateTeamMembership(membership)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist team membership inside the database", err}
			}
		}
	}
	return nil
}

func (handler *Handler) checkAdminPrivileges(user *portainer.User, userData *portainer.OAuthUserData, settings *portainer.Settings) *httperror.HandlerError {
	user.Role = portainer.StandardUserRole
	for _, team := range userData.Teams {
		if team == settings.OAuthSettings.AdminGroup {
			user.Role = portainer.AdministratorRole
		}
	}
	// Update user privileges
	err := handler.DataStore.User().UpdateUser(user.ID, user)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to update the user with necessary privileges", err}
	}
	return nil
}
