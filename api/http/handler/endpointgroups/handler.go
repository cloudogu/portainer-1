package endpointgroups

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/internal/authorization"

	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is the HTTP handler used to handle environment(endpoint) group operations.
type Handler struct {
	*mux.Router
	AuthorizationService *authorization.Service
	DataStore            dataservices.DataStore
}

// NewHandler creates a handler to manage environment(endpoint) group operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/endpoint_groups",
		bouncer.AdminAccess(httperror.LoggerHandler(h.endpointGroupCreate))).Methods(http.MethodPost)
	h.Handle("/endpoint_groups",
		bouncer.RestrictedAccess(httperror.LoggerHandler(h.endpointGroupList))).Methods(http.MethodGet)
	h.Handle("/endpoint_groups/{id}",
		bouncer.AdminAccess(httperror.LoggerHandler(h.endpointGroupInspect))).Methods(http.MethodGet)
	h.Handle("/endpoint_groups/{id}",
		bouncer.AdminAccess(httperror.LoggerHandler(h.endpointGroupUpdate))).Methods(http.MethodPut)
	h.Handle("/endpoint_groups/{id}",
		bouncer.AdminAccess(httperror.LoggerHandler(h.endpointGroupDelete))).Methods(http.MethodDelete)
	h.Handle("/endpoint_groups/{id}/endpoints/{endpointId}",
		bouncer.AdminAccess(httperror.LoggerHandler(h.endpointGroupAddEndpoint))).Methods(http.MethodPut)
	h.Handle("/endpoint_groups/{id}/endpoints/{endpointId}",
		bouncer.AdminAccess(httperror.LoggerHandler(h.endpointGroupDeleteEndpoint))).Methods(http.MethodDelete)
	return h
}
