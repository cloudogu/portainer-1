package middlewares

import (
	"net/http"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/gorilla/mux"
	httperror "github.com/portainer/libhttp/error"
)

func FeatureFlag(settingsService dataservices.SettingsService, feature portainer.Feature) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
			enabled := settingsService.IsFeatureFlagEnabled(feature)

			if !enabled {
				httperror.WriteError(rw, http.StatusForbidden, "This feature is not enabled", nil)
				return
			}

			next.ServeHTTP(rw, request)
		})
	}
}
