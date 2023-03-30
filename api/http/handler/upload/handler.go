package upload

import (
	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/http/security"
	httperror "github.com/portainer/libhttp/error"

	"net/http"

	"github.com/gorilla/mux"
)

// Handler is the HTTP handler used to handle upload operations.
type Handler struct {
	*mux.Router
	FileService portainer.FileService
}

// NewHandler creates a handler to manage upload operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/upload/tls/{certificate:(?:ca|cert|key)}",
		bouncer.AdminAccess(httperror.LoggerHandler(h.uploadTLS))).Methods(http.MethodPost)
	return h
}
