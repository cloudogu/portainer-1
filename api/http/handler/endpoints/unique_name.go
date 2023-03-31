package endpoints

import portainer "github.com/cloudogu/portainer-ce/api"

func (handler *Handler) isNameUnique(name string, endpointID portainer.EndpointID) (bool, error) {
	endpoints, err := handler.DataStore.Endpoint().Endpoints()
	if err != nil {
		return false, err
	}

	for _, endpoint := range endpoints {
		if endpoint.Name == name && (endpointID == 0 || endpoint.ID != endpointID) {
			return false, nil
		}
	}

	return true, nil
}
