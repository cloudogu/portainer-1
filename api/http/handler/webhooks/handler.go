package webhooks

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/dataservices"

	"github.com/cloudogu/portainer-ce/api/docker"
	"github.com/cloudogu/portainer-ce/api/http/security"
	httperror "github.com/portainer/libhttp/error"

	"github.com/gorilla/mux"
)

// Handler is the HTTP handler used to handle webhook operations.
type Handler struct {
	*mux.Router
	requestBouncer      *security.RequestBouncer
	DataStore           dataservices.DataStore
	DockerClientFactory *docker.ClientFactory
}

// NewHandler creates a handler to manage webhooks operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router:         mux.NewRouter(),
		requestBouncer: bouncer,
	}
	h.Handle("/webhooks",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.webhookCreate))).Methods(http.MethodPost)
	h.Handle("/webhooks/{id}",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.webhookUpdate))).Methods(http.MethodPut)
	h.Handle("/webhooks",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.webhookList))).Methods(http.MethodGet)
	h.Handle("/webhooks/{id}",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.webhookDelete))).Methods(http.MethodDelete)
	h.Handle("/webhooks/{token}",
		bouncer.PublicAccess(httperror.LoggerHandler(h.webhookExecute))).Methods(http.MethodPost)
	return h
}
