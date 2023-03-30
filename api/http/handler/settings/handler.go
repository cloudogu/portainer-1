package settings

import (
	"net/http"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/demo"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

func hideFields(settings *portainer.Settings) {
	settings.LDAPSettings.Password = ""
	settings.OAuthSettings.ClientSecret = ""
	settings.OAuthSettings.KubeSecretKey = nil
}

// Handler is the HTTP handler used to handle settings operations.
type Handler struct {
	*mux.Router
	DataStore       dataservices.DataStore
	FileService     portainer.FileService
	JWTService      dataservices.JWTService
	LDAPService     portainer.LDAPService
	SnapshotService portainer.SnapshotService
	demoService     *demo.Service
}

// NewHandler creates a handler to manage settings operations.
func NewHandler(bouncer *security.RequestBouncer, demoService *demo.Service) *Handler {
	h := &Handler{
		Router:      mux.NewRouter(),
		demoService: demoService,
	}
	h.Handle("/settings",
		bouncer.AdminAccess(httperror.LoggerHandler(h.settingsInspect))).Methods(http.MethodGet)
	h.Handle("/settings",
		bouncer.AdminAccess(httperror.LoggerHandler(h.settingsUpdate))).Methods(http.MethodPut)
	h.Handle("/settings/public",
		bouncer.PublicAccess(httperror.LoggerHandler(h.settingsPublic))).Methods(http.MethodGet)

	return h
}
