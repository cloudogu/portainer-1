package resourcecontrols

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is the HTTP handler used to handle resource control operations.
type Handler struct {
	*mux.Router
	DataStore dataservices.DataStore
}

// NewHandler creates a handler to manage resource control operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/resource_controls",
		bouncer.AdminAccess(httperror.LoggerHandler(h.resourceControlCreate))).Methods(http.MethodPost)
	h.Handle("/resource_controls/{id}",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.resourceControlUpdate))).Methods(http.MethodPut)
	h.Handle("/resource_controls/{id}",
		bouncer.AdminAccess(httperror.LoggerHandler(h.resourceControlDelete))).Methods(http.MethodDelete)
	return h
}
