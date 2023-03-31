package websocket

import (
	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/http/proxy/factory/kubernetes"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/cloudogu/portainer-ce/api/kubernetes/cli"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is the HTTP handler used to handle websocket operations.
type Handler struct {
	*mux.Router
	DataStore                   dataservices.DataStore
	SignatureService            portainer.DigitalSignatureService
	ReverseTunnelService        portainer.ReverseTunnelService
	KubernetesClientFactory     *cli.ClientFactory
	requestBouncer              *security.RequestBouncer
	connectionUpgrader          websocket.Upgrader
	kubernetesTokenCacheManager *kubernetes.TokenCacheManager
}

// NewHandler creates a handler to manage websocket operations.
func NewHandler(kubernetesTokenCacheManager *kubernetes.TokenCacheManager, bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router:                      mux.NewRouter(),
		connectionUpgrader:          websocket.Upgrader{},
		requestBouncer:              bouncer,
		kubernetesTokenCacheManager: kubernetesTokenCacheManager,
	}
	h.PathPrefix("/websocket/exec").Handler(
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.websocketExec)))
	h.PathPrefix("/websocket/attach").Handler(
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.websocketAttach)))
	h.PathPrefix("/websocket/pod").Handler(
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.websocketPodExec)))
	h.PathPrefix("/websocket/kubernetes-shell").Handler(
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.websocketShellPodExec)))
	return h
}
