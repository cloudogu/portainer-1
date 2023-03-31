package teams

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is the HTTP handler used to handle team operations.
type Handler struct {
	*mux.Router
	DataStore dataservices.DataStore
}

// NewHandler creates a handler to manage team operations.
func NewHandler(bouncer *security.RequestBouncer) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}

	adminRouter := h.NewRoute().Subrouter()
	adminRouter.Use(bouncer.AdminAccess)

	restrictedRouter := h.NewRoute().Subrouter()
	restrictedRouter.Use(bouncer.RestrictedAccess)

	teamLeaderRouter := h.NewRoute().Subrouter()
	teamLeaderRouter.Use(bouncer.TeamLeaderAccess)

	adminRouter.Handle("/teams", httperror.LoggerHandler(h.teamCreate)).Methods(http.MethodPost)
	restrictedRouter.Handle("/teams", httperror.LoggerHandler(h.teamList)).Methods(http.MethodGet)
	teamLeaderRouter.Handle("/teams/{id}", httperror.LoggerHandler(h.teamInspect)).Methods(http.MethodGet)
	adminRouter.Handle("/teams/{id}", httperror.LoggerHandler(h.teamUpdate)).Methods(http.MethodPut)
	adminRouter.Handle("/teams/{id}", httperror.LoggerHandler(h.teamDelete)).Methods(http.MethodDelete)
	teamLeaderRouter.Handle("/teams/{id}/memberships", httperror.LoggerHandler(h.teamMemberships)).Methods(http.MethodGet)

	return h
}
