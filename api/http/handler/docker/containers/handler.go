package containers

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/docker"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

type Handler struct {
	*mux.Router
	dockerClientFactory *docker.ClientFactory
}

// NewHandler creates a handler to process non-proxied requests to docker APIs directly.
func NewHandler(routePrefix string, bouncer *security.RequestBouncer, dockerClientFactory *docker.ClientFactory) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),

		dockerClientFactory: dockerClientFactory,
	}

	router := h.PathPrefix(routePrefix).Subrouter()
	router.Use(bouncer.AuthenticatedAccess)

	router.Handle("/{containerId}/gpus", httperror.LoggerHandler(h.containerGpusInspect)).Methods(http.MethodGet)

	return h
}
