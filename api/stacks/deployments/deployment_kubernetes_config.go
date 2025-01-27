package deployments

import (
	"fmt"
	"os"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/filesystem"
	k "github.com/cloudogu/portainer-ce/api/kubernetes"
	"github.com/cloudogu/portainer-ce/api/stacks/stackutils"
	"github.com/pkg/errors"
)

type KubernetesStackDeploymentConfig struct {
	stack             *portainer.Stack
	kuberneteDeployer portainer.KubernetesDeployer
	appLabels         k.KubeAppLabels
	user              *portainer.User
	endpoint          *portainer.Endpoint
	output            string
}

func CreateKubernetesStackDeploymentConfig(stack *portainer.Stack, kubeDeployer portainer.KubernetesDeployer, appLabels k.KubeAppLabels, user *portainer.User, endpoint *portainer.Endpoint) (*KubernetesStackDeploymentConfig, error) {

	return &KubernetesStackDeploymentConfig{
		stack:             stack,
		kuberneteDeployer: kubeDeployer,
		appLabels:         appLabels,
		user:              user,
		endpoint:          endpoint,
	}, nil
}

func (config *KubernetesStackDeploymentConfig) GetUsername() string {
	return config.user.Username
}

func (config *KubernetesStackDeploymentConfig) Deploy() error {
	fileNames := stackutils.GetStackFilePaths(config.stack, false)

	manifestFilePaths := make([]string, len(fileNames))

	tmpDir, err := os.MkdirTemp("", "kub_deployment")
	if err != nil {
		return errors.Wrap(err, "failed to create temp kub deployment directory")
	}

	defer os.RemoveAll(tmpDir)

	for _, fileName := range fileNames {
		manifestFilePath := filesystem.JoinPaths(tmpDir, fileName)
		manifestContent, err := os.ReadFile(filesystem.JoinPaths(config.stack.ProjectPath, fileName))
		if err != nil {
			return errors.Wrap(err, "failed to read manifest file")
		}

		if config.stack.IsComposeFormat {
			manifestContent, err = config.kuberneteDeployer.ConvertCompose(manifestContent)
			if err != nil {
				return errors.Wrap(err, "failed to convert docker compose file to a kube manifest")
			}
		}

		manifestContent, err = k.AddAppLabels(manifestContent, config.appLabels.ToMap())
		if err != nil {
			return errors.Wrap(err, "failed to add application labels")
		}

		err = filesystem.WriteToFile(manifestFilePath, []byte(manifestContent))
		if err != nil {
			return errors.Wrap(err, "failed to create temp manifest file")
		}

		manifestFilePaths = append(manifestFilePaths, manifestFilePath)
	}

	output, err := config.kuberneteDeployer.Deploy(config.user.ID, config.endpoint, manifestFilePaths, config.stack.Namespace)
	if err != nil {
		return fmt.Errorf("failed to deploy kubernete stack: %w", err)
	}

	config.output = output
	return nil
}

func (config *KubernetesStackDeploymentConfig) GetResponse() string {
	return config.output
}
