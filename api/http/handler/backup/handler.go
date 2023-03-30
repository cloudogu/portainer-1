package backup

import (
	"context"
	"net/http"

	"github.com/cloudogu/portainer-ce/api/adminmonitor"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/demo"
	"github.com/cloudogu/portainer-ce/api/http/middlewares"
	"github.com/cloudogu/portainer-ce/api/http/offlinegate"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is an http handler responsible for backup and restore portainer state
type Handler struct {
	*mux.Router
	bouncer         *security.RequestBouncer
	dataStore       dataservices.DataStore
	gate            *offlinegate.OfflineGate
	filestorePath   string
	shutdownTrigger context.CancelFunc
	adminMonitor    *adminmonitor.Monitor
}

// NewHandler creates an new instance of backup handler
func NewHandler(
	bouncer *security.RequestBouncer,
	dataStore dataservices.DataStore,
	gate *offlinegate.OfflineGate,
	filestorePath string,
	shutdownTrigger context.CancelFunc,
	adminMonitor *adminmonitor.Monitor,
	demoService *demo.Service,

) *Handler {

	h := &Handler{
		Router:          mux.NewRouter(),
		bouncer:         bouncer,
		dataStore:       dataStore,
		gate:            gate,
		filestorePath:   filestorePath,
		shutdownTrigger: shutdownTrigger,
		adminMonitor:    adminMonitor,
	}

	demoRestrictedRouter := h.NewRoute().Subrouter()
	demoRestrictedRouter.Use(middlewares.RestrictDemoEnv(demoService.IsDemo))

	demoRestrictedRouter.Handle("/backup", bouncer.RestrictedAccess(adminAccess(httperror.LoggerHandler(h.backup)))).Methods(http.MethodPost)
	demoRestrictedRouter.Handle("/restore", bouncer.PublicAccess(httperror.LoggerHandler(h.restore))).Methods(http.MethodPost)

	return h
}

func adminAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		securityContext, err := security.RetrieveRestrictedRequestContext(r)
		if err != nil {
			httperror.WriteError(w, http.StatusInternalServerError, "Unable to retrieve user info from request context", err)
		}

		if !securityContext.IsAdmin {
			httperror.WriteError(w, http.StatusUnauthorized, "User is not authorized to perform the action", nil)
		}

		next.ServeHTTP(w, r)
	})
}
