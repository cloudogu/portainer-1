package types

import portainer "github.com/cloudogu/portainer-ce/api"

type StoreManifestFunc func(stackFolder string, relatedEndpointIds []portainer.EndpointID) (composePath, manifestPath, projectPath string, err error)
