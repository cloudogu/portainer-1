package helm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/datastore"
	"github.com/cloudogu/portainer-ce/api/exec/exectest"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/cloudogu/portainer-ce/api/jwt"
	"github.com/cloudogu/portainer-ce/api/kubernetes"
	"github.com/portainer/portainer/pkg/libhelm/binary/test"
	"github.com/portainer/portainer/pkg/libhelm/options"
	"github.com/stretchr/testify/assert"

	helper "github.com/cloudogu/portainer-ce/api/internal/testhelpers"
)

func Test_helmDelete(t *testing.T) {
	is := assert.New(t)

	_, store, teardown := datastore.MustNewTestStore(t, true, true)
	defer teardown()

	err := store.Endpoint().Create(&portainer.Endpoint{ID: 1})
	is.NoError(err, "Error creating environment")

	err = store.User().Create(&portainer.User{Username: "admin", Role: portainer.AdministratorRole})
	is.NoError(err, "Error creating a user")

	jwtService, err := jwt.NewService("1h", store)
	is.NoError(err, "Error initiating jwt service")

	kubernetesDeployer := exectest.NewKubernetesDeployer()
	helmPackageManager := test.NewMockHelmBinaryPackageManager("")
	kubeClusterAccessService := kubernetes.NewKubeClusterAccessService("", "", "")
	h := NewHandler(helper.NewTestRequestBouncer(), store, jwtService, kubernetesDeployer, helmPackageManager, kubeClusterAccessService)

	is.NotNil(h, "Handler should not fail")

	// Install a single chart directly, to be deleted by the handler
	options := options.InstallOptions{Name: "nginx-1", Chart: "nginx", Namespace: "default"}
	h.helmPackageManager.Install(options)

	t.Run("helmDelete succeeds with admin user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/1/kubernetes/helm/%s", options.Name), nil)
		ctx := security.StoreTokenData(req, &portainer.TokenData{ID: 1, Username: "admin", Role: 1})
		req = req.WithContext(ctx)
		req.Header.Add("Authorization", "Bearer dummytoken")

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)

		is.Equal(http.StatusNoContent, rr.Code, "Status should be 204")
	})
}
