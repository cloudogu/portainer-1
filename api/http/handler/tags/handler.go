package tags

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is the HTTP handler used to handle tag operations.
type Handler struct {
	*mux.Router
	DataStore dataservices.DataStore
}

// NewHandler creates a handler to manage tag operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/tags",
		bouncer.AdminAccess(httperror.LoggerHandler(h.tagCreate))).Methods(http.MethodPost)
	h.Handle("/tags",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.tagList))).Methods(http.MethodGet)
	h.Handle("/tags/{id}",
		bouncer.AdminAccess(httperror.LoggerHandler(h.tagDelete))).Methods(http.MethodDelete)

	return h
}
