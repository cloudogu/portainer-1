package kubernetes

import (
	"context"
	"time"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/kubernetes/cli"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Snapshotter struct {
	clientFactory *cli.ClientFactory
}

// NewSnapshotter returns a new Snapshotter instance
func NewSnapshotter(clientFactory *cli.ClientFactory) *Snapshotter {
	return &Snapshotter{
		clientFactory: clientFactory,
	}
}

// CreateSnapshot creates a snapshot of a specific Kubernetes environment(endpoint)
func (snapshotter *Snapshotter) CreateSnapshot(endpoint *portainer.Endpoint) (*portainer.KubernetesSnapshot, error) {
	client, err := snapshotter.clientFactory.CreateClient(endpoint)
	if err != nil {
		return nil, err
	}

	return snapshot(client, endpoint)
}

func snapshot(cli *kubernetes.Clientset, endpoint *portainer.Endpoint) (*portainer.KubernetesSnapshot, error) {
	res := cli.RESTClient().Get().AbsPath("/healthz").Do(context.TODO())
	if res.Error() != nil {
		return nil, res.Error()
	}

	snapshot := &portainer.KubernetesSnapshot{}

	err := snapshotVersion(snapshot, cli)
	if err != nil {
		log.Warn().Str("endpoint", endpoint.Name).Err(err).Msg("unable to snapshot cluster version")
	}

	err = snapshotNodes(snapshot, cli)
	if err != nil {
		log.Warn().Str("endpoint", endpoint.Name).Err(err).Msg("unable to snapshot cluster nodes")
	}

	snapshot.Time = time.Now().Unix()
	return snapshot, nil
}

func snapshotVersion(snapshot *portainer.KubernetesSnapshot, cli *kubernetes.Clientset) error {
	versionInfo, err := cli.ServerVersion()
	if err != nil {
		return err
	}

	snapshot.KubernetesVersion = versionInfo.GitVersion
	return nil
}

func snapshotNodes(snapshot *portainer.KubernetesSnapshot, cli *kubernetes.Clientset) error {
	nodeList, err := cli.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	var totalCPUs, totalMemory int64
	for _, node := range nodeList.Items {
		totalCPUs += node.Status.Capacity.Cpu().Value()
		totalMemory += node.Status.Capacity.Memory().Value()
	}

	snapshot.TotalCPU = totalCPUs
	snapshot.TotalMemory = totalMemory
	snapshot.NodeCount = len(nodeList.Items)
	return nil
}
