package websocket

import (
	"github.com/cloudogu/portainer-ce/api"
)

type webSocketRequestParams struct {
	ID       string
	nodeName string
	endpoint *portainer.Endpoint
}
