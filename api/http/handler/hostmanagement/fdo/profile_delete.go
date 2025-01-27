package fdo

import (
	"errors"
	"net/http"

	portainer "github.com/cloudogu/portainer-ce/api"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

// @id deleteProfile
// @summary deletes a FDO Profile
// @description deletes a FDO Profile
// @description **Access policy**: administrator
// @tags intel
// @security jwt
// @produce json
// @success 200 "Success"
// @failure 400 "Invalid request"
// @failure 500 "Server error"
// @router /fdo/profiles/{id} [delete]
func (handler *Handler) deleteProfile(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	id, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Bad request", errors.New("missing 'id' query parameter"))
	}

	err = handler.DataStore.FDOProfile().Delete(portainer.FDOProfileID(id))
	if err != nil {
		return httperror.InternalServerError("Unable to delete Profile", err)
	}

	return response.Empty(w)
}
