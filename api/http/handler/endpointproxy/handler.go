package endpointproxy

import (
	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/http/proxy"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is the HTTP handler used to proxy requests to external APIs.
type Handler struct {
	*mux.Router
	DataStore            dataservices.DataStore
	requestBouncer       *security.RequestBouncer
	ProxyManager         *proxy.Manager
	ReverseTunnelService portainer.ReverseTunnelService
}

// NewHandler creates a handler to proxy requests to external APIs.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router:         mux.NewRouter(),
		requestBouncer: bouncer,
	}
	h.PathPrefix("/{id}/azure").Handler(
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.proxyRequestsToAzureAPI)))
	h.PathPrefix("/{id}/docker").Handler(
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.proxyRequestsToDockerAPI)))
	h.PathPrefix("/{id}/kubernetes").Handler(
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.proxyRequestsToKubernetesAPI)))
	h.PathPrefix("/{id}/agent/docker").Handler(
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.proxyRequestsToDockerAPI)))
	h.PathPrefix("/{id}/agent/kubernetes").Handler(
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.proxyRequestsToKubernetesAPI)))
	return h
}
