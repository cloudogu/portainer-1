package roles

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is the HTTP handler used to handle role operations.
type Handler struct {
	*mux.Router
	DataStore dataservices.DataStore
}

// NewHandler creates a handler to manage role operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/roles",
		bouncer.AdminAccess(httperror.LoggerHandler(h.roleList))).Methods(http.MethodGet)

	return h
}
