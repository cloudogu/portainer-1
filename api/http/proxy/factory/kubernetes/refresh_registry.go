package kubernetes

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/internal/registryutils"
)

func (transport *baseTransport) refreshRegistry(request *http.Request, namespace string) (err error) {
	cli, err := transport.k8sClientFactory.GetKubeClient(transport.endpoint)
	if err != nil {
		return
	}

	err = registryutils.RefreshEcrSecret(cli, transport.endpoint, transport.dataStore, namespace)

	return
}
