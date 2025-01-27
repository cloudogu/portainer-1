package users

import (
	"net/http"

	portainer "github.com/cloudogu/portainer-ce/api"
	httperrors "github.com/cloudogu/portainer-ce/api/http/errors"
	"github.com/cloudogu/portainer-ce/api/http/security"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

// @id UserGetAPIKeys
// @summary Get all API keys for a user
// @description Gets all API keys for a user.
// @description Only the calling user or admin can retrieve api-keys.
// @description **Access policy**: authenticated
// @tags users
// @security ApiKeyAuth
// @security jwt
// @produce json
// @param id path int true "User identifier"
// @success 200 {array} portainer.APIKey "Success"
// @failure 400 "Invalid request"
// @failure 403 "Permission denied"
// @failure 404 "User not found"
// @failure 500 "Server error"
// @router /users/{id}/tokens [get]
func (handler *Handler) userGetAccessTokens(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	userID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid user identifier route variable", err)
	}

	tokenData, err := security.RetrieveTokenData(r)
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve user authentication token", err)
	}

	if tokenData.Role != portainer.AdministratorRole && tokenData.ID != portainer.UserID(userID) {
		return httperror.Forbidden("Permission denied to get user access tokens", httperrors.ErrUnauthorized)
	}

	_, err = handler.DataStore.User().User(portainer.UserID(userID))
	if err != nil {
		if handler.DataStore.IsErrObjectNotFound(err) {
			return httperror.NotFound("Unable to find a user with the specified identifier inside the database", err)
		}
		return httperror.InternalServerError("Unable to find a user with the specified identifier inside the database", err)
	}

	apiKeys, err := handler.apiKeyService.GetAPIKeys(portainer.UserID(userID))
	if err != nil {
		return httperror.InternalServerError("Internal Server Error", err)
	}

	for idx := range apiKeys {
		hideAPIKeyFields(&apiKeys[idx])
	}

	return response.JSON(w, apiKeys)
}

// hideAPIKeyFields remove the digest from the API key (it is not needed in the response)
func hideAPIKeyFields(apiKey *portainer.APIKey) {
	apiKey.Digest = nil
}
