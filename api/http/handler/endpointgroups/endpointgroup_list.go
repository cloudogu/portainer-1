package endpointgroups

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/http/security"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
)

// GET request on /api/endpoint_groups
func (handler *Handler) endpointGroupList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	endpointGroups, err := handler.DataStore.EndpointGroup().EndpointGroups()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve endpoint groups from the database", err}
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve info from request context", err}
	}

	endpointGroups = security.FilterEndpointGroups(endpointGroups, securityContext)
	return response.JSON(w, endpointGroups)
}
