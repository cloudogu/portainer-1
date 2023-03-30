package motd

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/gorilla/mux"
)

// Handler is the HTTP handler used to handle MOTD operations.
type Handler struct {
	*mux.Router
}

// NewHandler returns a new Handler
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}
	h.Handle("/motd",
		bouncer.RestrictedAccess(http.HandlerFunc(h.motd))).Methods(http.MethodGet)

	return h
}
