package fdo

import (
	"net/http"

	"github.com/gorilla/mux"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/http/security"
	httperror "github.com/portainer/libhttp/error"
)

type Handler struct {
	*mux.Router
	DataStore   dataservices.DataStore
	FileService portainer.FileService
}

func NewHandler(bouncer *security.RequestBouncer, dataStore dataservices.DataStore, fileService portainer.FileService) *Handler {
	h := &Handler{
		Router:      mux.NewRouter(),
		DataStore:   dataStore,
		FileService: fileService,
	}

	h.Handle("/fdo/configure", bouncer.AdminAccess(httperror.LoggerHandler(h.fdoConfigure))).Methods(http.MethodPost)
	h.Handle("/fdo/list", bouncer.AdminAccess(httperror.LoggerHandler(h.fdoListAll))).Methods(http.MethodGet)
	h.Handle("/fdo/register", bouncer.AdminAccess(httperror.LoggerHandler(h.fdoRegisterDevice))).Methods(http.MethodPost)
	h.Handle("/fdo/configure/{guid}", bouncer.AdminAccess(httperror.LoggerHandler(h.fdoConfigureDevice))).Methods(http.MethodPost)

	h.Handle("/fdo/profiles", bouncer.AdminAccess(httperror.LoggerHandler(h.fdoProfileList))).Methods(http.MethodGet)
	h.Handle("/fdo/profiles", bouncer.AdminAccess(httperror.LoggerHandler(h.createProfile))).Methods(http.MethodPost)
	h.Handle("/fdo/profiles/{id}", bouncer.AdminAccess(httperror.LoggerHandler(h.fdoProfileInspect))).Methods(http.MethodGet)
	h.Handle("/fdo/profiles/{id}", bouncer.AdminAccess(httperror.LoggerHandler(h.updateProfile))).Methods(http.MethodPut)
	h.Handle("/fdo/profiles/{id}", bouncer.AdminAccess(httperror.LoggerHandler(h.deleteProfile))).Methods(http.MethodDelete)
	h.Handle("/fdo/profiles/{id}/duplicate", bouncer.AdminAccess(httperror.LoggerHandler(h.duplicateProfile))).Methods(http.MethodPost)

	return h
}
