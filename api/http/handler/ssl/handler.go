package ssl

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/cloudogu/portainer-ce/api/internal/ssl"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is the HTTP handler used to handle MOTD operations.
type Handler struct {
	*mux.Router
	SSLService *ssl.Service
}

// NewHandler returns a new Handler
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/ssl",
		bouncer.AdminAccess(httperror.LoggerHandler(h.sslInspect))).Methods(http.MethodGet)
	h.Handle("/ssl",
		bouncer.AdminAccess(httperror.LoggerHandler(h.sslUpdate))).Methods(http.MethodPut)

	return h
}
