package edgestacks

import (
	"fmt"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/internal/endpointutils"
)

func hasKubeEndpoint(endpointService dataservices.EndpointService, endpointIDs []portainer.EndpointID) (bool, error) {
	return hasEndpointPredicate(endpointService, endpointIDs, endpointutils.IsKubernetesEndpoint)
}

func hasDockerEndpoint(endpointService dataservices.EndpointService, endpointIDs []portainer.EndpointID) (bool, error) {
	return hasEndpointPredicate(endpointService, endpointIDs, endpointutils.IsDockerEndpoint)
}

func hasEndpointPredicate(endpointService dataservices.EndpointService, endpointIDs []portainer.EndpointID, predicate func(*portainer.Endpoint) bool) (bool, error) {
	for _, endpointID := range endpointIDs {
		endpoint, err := endpointService.Endpoint(endpointID)
		if err != nil {
			return false, fmt.Errorf("failed to retrieve environment from database: %w", err)
		}

		if predicate(endpoint) {
			return true, nil
		}
	}

	return false, nil
}
