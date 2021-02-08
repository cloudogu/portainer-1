package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/portainer/api"
	bolterrors "github.com/portainer/portainer/api/bolt/errors"
	httperrors "github.com/portainer/portainer/api/http/errors"
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

	user, err := handler.DataStore.User().UserByUsername(userData.Username)
	if err != nil && err != bolterrors.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve a user with the specified username from the database", err}
	}

	if user == nil && !settings.OAuthSettings.OAuthAutoCreateUsers {
		return &httperror.HandlerError{http.StatusForbidden, "Account not created beforehand in Portainer and automatic user provisioning not enabled", httperrors.ErrUnauthorized}
	}

	if user == nil {
		user = &portainer.User{
			Username: userData.Username,
			Role:     portainer.StandardUserRole,
		}

		err = handler.DataStore.User().CreateUser(user)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist user inside the database", err}
		}

		if settings.OAuthSettings.DefaultTeamID != 0 {
			membership := &portainer.TeamMembership{
				UserID: user.ID,
				TeamID: settings.OAuthSettings.DefaultTeamID,
				Role:   portainer.TeamMember,
			}

			err = handler.DataStore.TeamMembership().CreateTeamMembership(membership)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist team membership inside the database", err}
			}
		}
	}
	user.OAuthToken = userData.OAuthToken

	return handler.writeToken(w, user)
}

func (handler *Handler) userProvisioning(userData *portainer.OAuthUserData, settings *portainer.Settings) *httperror.HandlerError {
	user, err := handler.DataStore.User().UserByUsername(userData.Username)
	if err != nil && err != bolterrors.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve a user with the specified username from the database", err}
	}

	if user == nil && !settings.OAuthSettings.OAuthAutoCreateUsers {
		return &httperror.HandlerError{http.StatusForbidden, "Account not created beforehand in Portainer and automatic user provisioning not enabled", httperrors.ErrUnauthorized}
	}

	if user == nil {
		user = &portainer.User{
			Username: userData.Username,
			Role:     portainer.StandardUserRole,
		}

		err = handler.DataStore.User().CreateUser(user)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist user inside the database", err}
		}

		if settings.OAuthSettings.DefaultTeamID != 0 {
			membership := &portainer.TeamMembership{
				UserID: user.ID,
				TeamID: settings.OAuthSettings.DefaultTeamID,
				Role:   portainer.TeamMember,
			}

			err = handler.DataStore.TeamMembership().CreateTeamMembership(membership)
			if err != nil {
				return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist team membership inside the database", err}
			}
		}
	}
	user.OAuthToken = userData.OAuthToken

	return nil
}
