package auth

import (
	"net/http"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/http/proxy"
	"github.com/cloudogu/portainer-ce/api/http/proxy/factory/kubernetes"
	"github.com/cloudogu/portainer-ce/api/http/security"
	httperror "github.com/portainer/libhttp/error"

	"github.com/gorilla/mux"
)

// Handler is the HTTP handler used to handle authentication operations.
type Handler struct {
	*mux.Router
	DataStore                   dataservices.DataStore
	CryptoService               portainer.CryptoService
	JWTService                  dataservices.JWTService
	LDAPService                 portainer.LDAPService
	OAuthService                portainer.OAuthService
	ProxyManager                *proxy.Manager
	KubernetesTokenCacheManager *kubernetes.TokenCacheManager
	passwordStrengthChecker     security.PasswordStrengthChecker
}

// NewHandler creates a handler to manage authentication operations.
func NewHandler(bouncer *security.RequestBouncer, rateLimiter *security.RateLimiter, passwordStrengthChecker security.PasswordStrengthChecker) *Handler {
	h := &Handler{
		Router:                  mux.NewRouter(),
		passwordStrengthChecker: passwordStrengthChecker,
	}

	h.Handle("/auth/oauth/validate",
		rateLimiter.LimitAccess(bouncer.PublicAccess(httperror.LoggerHandler(h.validateOAuth)))).Methods(http.MethodPost)
	h.Handle("/auth/oauth/logout",
		rateLimiter.LimitAccess(bouncer.PublicAccess(httperror.LoggerHandler(h.invalidateOAuthSession)))).Methods(http.MethodPost)
	h.Handle("/auth/oauth/verifyToken",
		rateLimiter.LimitAccess(bouncer.PublicAccess(httperror.LoggerHandler(h.isJWTTokenNotBlocked)))).Methods(http.MethodPost)
	h.Handle("/auth/oauth/apiToken",
		rateLimiter.LimitAccess(bouncer.PublicAccess(httperror.LoggerHandler(h.authenticateViaApi)))).Methods(http.MethodPost)
	h.Handle("/auth",
		rateLimiter.LimitAccess(bouncer.PublicAccess(httperror.LoggerHandler(h.authenticate)))).Methods(http.MethodPost)
	h.Handle("/auth/logout",
		bouncer.AuthenticatedAccess(httperror.LoggerHandler(h.logout))).Methods(http.MethodPost)

	return h
}
