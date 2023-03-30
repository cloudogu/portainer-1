package openamt

import (
	"net/http"

	"github.com/gorilla/mux"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/docker"
	"github.com/cloudogu/portainer-ce/api/http/security"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is the HTTP handler used to handle OpenAMT operations.
type Handler struct {
	*mux.Router
	OpenAMTService      portainer.OpenAMTService
	DataStore           dataservices.DataStore
	DockerClientFactory *docker.ClientFactory
}

// NewHandler returns a new Handler
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}

	h.Handle("/open_amt/configure", bouncer.AdminAccess(httperror.LoggerHandler(h.openAMTConfigure))).Methods(http.MethodPost)
	h.Handle("/open_amt/{id}/info", bouncer.AdminAccess(httperror.LoggerHandler(h.openAMTHostInfo))).Methods(http.MethodGet)
	h.Handle("/open_amt/{id}/activate", bouncer.AdminAccess(httperror.LoggerHandler(h.openAMTActivate))).Methods(http.MethodPost)
	h.Handle("/open_amt/{id}/devices", bouncer.AdminAccess(httperror.LoggerHandler(h.openAMTDevices))).Methods(http.MethodGet)
	h.Handle("/open_amt/{id}/devices/{deviceId}/action", bouncer.AdminAccess(httperror.LoggerHandler(h.deviceAction))).Methods(http.MethodPost)
	h.Handle("/open_amt/{id}/devices/{deviceId}/features", bouncer.AdminAccess(httperror.LoggerHandler(h.deviceFeatures))).Methods(http.MethodPost)

	return h
}
