package users

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/bolt/errors"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
)

// GET request on /api/users/admin/check
func (handler *Handler) adminCheck(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	users, err := handler.DataStore.User().UsersByRole(portainer.AdministratorRole)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve users from the database", err}
	}

	if len(users) == 0 {
		return &httperror.HandlerError{http.StatusNotFound, "No administrator account found inside the database", errors.ErrObjectNotFound}
	}

	return response.Empty(w)
}
