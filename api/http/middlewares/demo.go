package middlewares

import (
	"net/http"

	"github.com/cloudogu/portainer-ce/api/http/errors"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

// restrict functionality on demo environments
func RestrictDemoEnv(isDemo func() bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isDemo() {
				next.ServeHTTP(w, r)
				return
			}

			httperror.WriteError(w, http.StatusBadRequest, errors.ErrNotAvailableInDemo.Error(), errors.ErrNotAvailableInDemo)
		})
	}
}
