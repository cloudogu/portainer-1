package exec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/internal/testhelpers"
	"github.com/portainer/docker-compose-wrapper/compose"

	"github.com/rs/zerolog/log"
)

const composeFile = `version: "3.9"
services:
  busybox:
    image: "alpine:latest"
    container_name: "compose_wrapper_test"`
const composedContainerName = "compose_wrapper_test"

func setup(t *testing.T) (*portainer.Stack, *portainer.Endpoint) {
	dir := t.TempDir()
	composeFileName := "compose_wrapper_test.yml"
	f, _ := os.Create(filepath.Join(dir, composeFileName))
	f.WriteString(composeFile)

	stack := &portainer.Stack{
		ProjectPath: dir,
		EntryPoint:  composeFileName,
		Name:        "project-name",
	}

	endpoint := &portainer.Endpoint{
		URL: "unix://",
	}

	return stack, endpoint
}

func Test_UpAndDown(t *testing.T) {

	testhelpers.IntegrationTest(t)

	stack, endpoint := setup(t)

	deployer, err := compose.NewComposeDeployer("", "")
	if err != nil {
		t.Fatal(err)
	}

	w, err := NewComposeStackManager(deployer, nil)
	if err != nil {
		t.Fatalf("Failed creating manager: %s", err)
	}

	ctx := context.TODO()

	err = w.Up(ctx, stack, endpoint, false)
	if err != nil {
		t.Fatalf("Error calling docker-compose up: %s", err)
	}

	if !containerExists(composedContainerName) {
		t.Fatal("container should exist")
	}

	err = w.Down(ctx, stack, endpoint)
	if err != nil {
		t.Fatalf("Error calling docker-compose down: %s", err)
	}

	if containerExists(composedContainerName) {
		t.Fatal("container should be removed")
	}
}

func containerExists(containerName string) bool {
	cmd := exec.Command("docker", "ps", "-a", "-f", fmt.Sprintf("name=%s", containerName))

	out, err := cmd.Output()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list containers")
	}

	return strings.Contains(string(out), containerName)
}
