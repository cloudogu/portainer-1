package auth

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/http/proxy"
	"github.com/cloudogu/portainer-ce/api/http/proxy/factory/kubernetes"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

// Handler is the HTTP handler used to handle authentication operations.
type Handler struct {
	*mux.Router
	DataStore                   portainer.DataStore
	CryptoService               portainer.CryptoService
	JWTService                  portainer.JWTService
	LDAPService                 portainer.LDAPService
	OAuthService                portainer.OAuthService
	ProxyManager                *proxy.Manager
	KubernetesTokenCacheManager *kubernetes.TokenCacheManager
}

// NewHandler creates a handler to manage authentication operations.
func NewHandler(bouncer *security.RequestBouncer, rateLimiter *security.RateLimiter) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
	}

	h.Handle("/auth/oauth/validate",
		rateLimiter.LimitAccess(bouncer.PublicAccess(httperror.LoggerHandler(h.validateOAuth)))).Methods(http.MethodPost)
	h.Handle("/auth/oauth/logout",
		rateLimiter.LimitAccess(bouncer.PublicAccess(httperror.LoggerHandler(h.invalidateOAuthSession)))).Methods(http.MethodPost)
	h.Handle("/auth/oauth/verifyToken",
		rateLimiter.LimitAccess(bouncer.PublicAccess(httperror.LoggerHandler(h.isJWTTokenNotBlocked)))).Methods(http.MethodPost)
	h.Handle("/auth",
		rateLimiter.LimitAccess(bouncer.PublicAccess(httperror.LoggerHandler(h.authenticate)))).Methods(http.MethodPost)
	h.Handle("/auth/logout",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.logout))).Methods(http.MethodPost)

	return h
}
