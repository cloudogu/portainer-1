package factory

import (
	"net/http"
	"net/url"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/http/proxy/factory/azure"
)

func newAzureProxy(endpoint *portainer.Endpoint) (http.Handler, error) {
	remoteURL, err := url.Parse(azureAPIBaseURL)
	if err != nil {
		return nil, err
	}

	proxy := newSingleHostReverseProxyWithHostHeader(remoteURL)
	proxy.Transport = azure.NewTransport(&endpoint.AzureCredentials)
	return proxy, nil
}
