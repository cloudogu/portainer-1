package websocket

import portainer "github.com/cloudogu/portainer-ce/api"

type webSocketRequestParams struct {
	ID       string
	nodeName string
	endpoint *portainer.Endpoint
	token    string
}
