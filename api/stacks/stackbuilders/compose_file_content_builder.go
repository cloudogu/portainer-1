package stackbuilders

import (
	"strconv"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/filesystem"
	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/cloudogu/portainer-ce/api/stacks/deployments"
	httperror "github.com/portainer/libhttp/error"
)

type ComposeStackFileContentBuilder struct {
	FileContentMethodStackBuilder
	SecurityContext *security.RestrictedRequestContext
}

// CreateComposeStackFileContentBuilder creates a builder for the compose stack (docker standalone) that will be deployed by file content method
func CreateComposeStackFileContentBuilder(securityContext *security.RestrictedRequestContext,
	dataStore dataservices.DataStore,
	fileService portainer.FileService,
	stackDeployer deployments.StackDeployer) *ComposeStackFileContentBuilder {

	return &ComposeStackFileContentBuilder{
		FileContentMethodStackBuilder: FileContentMethodStackBuilder{
			StackBuilder: CreateStackBuilder(dataStore, fileService, stackDeployer),
		},
		SecurityContext: securityContext,
	}
}

func (b *ComposeStackFileContentBuilder) SetGeneralInfo(payload *StackPayload, endpoint *portainer.Endpoint) FileContentMethodStackBuildProcess {
	b.FileContentMethodStackBuilder.SetGeneralInfo(payload, endpoint)
	return b
}

func (b *ComposeStackFileContentBuilder) SetUniqueInfo(payload *StackPayload) FileContentMethodStackBuildProcess {
	if b.hasError() {
		return b
	}
	b.stack.Name = payload.Name
	b.stack.Type = portainer.DockerComposeStack
	b.stack.EntryPoint = filesystem.ComposeFileDefaultName
	b.stack.Env = payload.Env
	b.stack.FromAppTemplate = payload.FromAppTemplate
	return b
}

func (b *ComposeStackFileContentBuilder) SetFileContent(payload *StackPayload) FileContentMethodStackBuildProcess {
	if b.hasError() {
		return b
	}

	stackFolder := strconv.Itoa(int(b.stack.ID))
	projectPath, err := b.fileService.StoreStackFileFromBytes(stackFolder, b.stack.EntryPoint, []byte(payload.StackFileContent))
	if err != nil {
		b.err = httperror.InternalServerError("Unable to persist Compose file on disk", err)
		return b
	}
	b.stack.ProjectPath = projectPath

	return b
}

func (b *ComposeStackFileContentBuilder) Deploy(payload *StackPayload, endpoint *portainer.Endpoint) FileContentMethodStackBuildProcess {
	if b.hasError() {
		return b
	}

	composeDeploymentConfig, err := deployments.CreateComposeStackDeploymentConfig(b.SecurityContext, b.stack, endpoint, b.dataStore, b.fileService, b.stackDeployer, false, false)
	if err != nil {
		b.err = httperror.InternalServerError(err.Error(), err)
		return b
	}

	b.deploymentConfiger = composeDeploymentConfig
	b.stack.CreatedBy = b.deploymentConfiger.GetUsername()

	return b.FileContentMethodStackBuilder.Deploy(payload, endpoint)
}
