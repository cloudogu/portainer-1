package docker

import (
	"errors"
	"net/http"

	"github.com/cloudogu/portainer-ce/api/docker"
	"github.com/cloudogu/portainer-ce/api/internal/endpointutils"

	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/http/handler/docker/containers"
	"github.com/cloudogu/portainer-ce/api/http/middlewares"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/cloudogu/portainer-ce/api/internal/authorization"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is the HTTP handler which will natively deal with to external environments(endpoints).
type Handler struct {
	*mux.Router
	requestBouncer       *security.RequestBouncer
	dataStore            dataservices.DataStore
	dockerClientFactory  *docker.ClientFactory
	authorizationService *authorization.Service
}

// NewHandler creates a handler to process non-proxied requests to docker APIs directly.
func NewHandler(bouncer *security.RequestBouncer, authorizationService *authorization.Service, dataStore dataservices.DataStore, dockerClientFactory *docker.ClientFactory) *Handler {
	h := &Handler{
		Router:               mux.NewRouter(),
		requestBouncer:       bouncer,
		authorizationService: authorizationService,
		dataStore:            dataStore,
		dockerClientFactory:  dockerClientFactory,
	}

	// endpoints
	endpointRouter := h.PathPrefix("/{id}").Subrouter()
	endpointRouter.Use(middlewares.WithEndpoint(dataStore.Endpoint(), "id"))
	endpointRouter.Use(dockerOnlyMiddleware)

	containersHandler := containers.NewHandler("/{id}/containers", bouncer, dockerClientFactory)
	endpointRouter.PathPrefix("/containers").Handler(containersHandler)
	return h
}

func dockerOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		endpoint, err := middlewares.FetchEndpoint(request)
		if err != nil {
			httperror.WriteError(rw, http.StatusInternalServerError, "Unable to find an environment on request context", err)
			return
		}

		if !endpointutils.IsDockerEndpoint(endpoint) {
			errMessage := "environment is not a docker environment"
			httperror.WriteError(rw, http.StatusBadRequest, errMessage, errors.New(errMessage))
			return
		}
		next.ServeHTTP(rw, request)
	})
}
