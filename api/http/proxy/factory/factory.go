package factory

import (
	"net/http"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/http/proxy/factory/kubernetes"

	"github.com/cloudogu/portainer-ce/api/kubernetes/cli"

	"github.com/cloudogu/portainer-ce/api/docker"
)

const azureAPIBaseURL = "https://management.azure.com"

type (
	// ProxyFactory is a factory to create reverse proxies
	ProxyFactory struct {
		dataStore                   dataservices.DataStore
		signatureService            portainer.DigitalSignatureService
		reverseTunnelService        portainer.ReverseTunnelService
		dockerClientFactory         *docker.ClientFactory
		kubernetesClientFactory     *cli.ClientFactory
		kubernetesTokenCacheManager *kubernetes.TokenCacheManager
		gitService                  portainer.GitService
	}
)

// NewProxyFactory returns a pointer to a new instance of a ProxyFactory
func NewProxyFactory(dataStore dataservices.DataStore, signatureService portainer.DigitalSignatureService, tunnelService portainer.ReverseTunnelService, clientFactory *docker.ClientFactory, kubernetesClientFactory *cli.ClientFactory, kubernetesTokenCacheManager *kubernetes.TokenCacheManager, gitService portainer.GitService) *ProxyFactory {
	return &ProxyFactory{
		dataStore:                   dataStore,
		signatureService:            signatureService,
		reverseTunnelService:        tunnelService,
		dockerClientFactory:         clientFactory,
		kubernetesClientFactory:     kubernetesClientFactory,
		kubernetesTokenCacheManager: kubernetesTokenCacheManager,
		gitService:                  gitService,
	}
}

// NewEndpointProxy returns a new reverse proxy (filesystem based or HTTP) to an environment(endpoint) API server
func (factory *ProxyFactory) NewEndpointProxy(endpoint *portainer.Endpoint) (http.Handler, error) {
	switch endpoint.Type {
	case portainer.AzureEnvironment:
		return newAzureProxy(endpoint, factory.dataStore)
	case portainer.EdgeAgentOnKubernetesEnvironment, portainer.AgentOnKubernetesEnvironment, portainer.KubernetesLocalEnvironment:
		return factory.newKubernetesProxy(endpoint)
	}

	return factory.newDockerProxy(endpoint)
}

// NewGitlabProxy returns a new HTTP proxy to a Gitlab API server
func (factory *ProxyFactory) NewGitlabProxy(gitlabAPIUri string) (http.Handler, error) {
	return newGitlabProxy(gitlabAPIUri)
}
