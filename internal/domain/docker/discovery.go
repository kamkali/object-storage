package docker

import (
	"fmt"
	"net"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/kamkalis/object-storage/internal/domain/minio"
	"golang.org/x/net/context"
)

// Discoverer can discover storage nodes and map them to concrete nodes
type Discoverer struct {
	nodeIdentifier    string
	networkIdentifier string
	nodePort          string
	c                 *client.Client
}

const (
	ImageMinio = "minio/minio"

	EnvKeyMinioAccessKey = "MINIO_ACCESS_KEY"
	EnvKeyMinioSecretKey = "MINIO_SECRET_KEY"
)

func NewNodeDiscoverer(nodeID, nodePort, networkID string) (*Discoverer, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("new docker client: %w", err)
	}

	return &Discoverer{
		c:                 cli,
		nodeIdentifier:    nodeID,
		nodePort:          nodePort,
		networkIdentifier: networkID,
	}, nil
}

func (n *Discoverer) DiscoverNodes(ctx context.Context) ([]domain.StorageNode, error) {
	containers, err := n.c.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("name", n.nodeIdentifier)),
	})
	if err != nil {
		return nil, fmt.Errorf("list docker containers: %w", err)
	}

	var domainNodes []domain.StorageNode
	for _, container := range containers {
		node, err := n.toDomainNode(ctx, container)
		if err != nil {
			return nil, fmt.Errorf("maping docker node to domain: %w", err)
		}
		domainNodes = append(domainNodes, node)
	}

	return domainNodes, nil
}

func (n *Discoverer) toDomainNode(ctx context.Context, c types.Container) (domain.StorageNode, error) {
	switch c.Image {
	case ImageMinio:
		return n.mapToMinioNode(ctx, c)
	default:
		return nil, fmt.Errorf("unknown storage image")
	}
}

func (n *Discoverer) mapToMinioNode(ctx context.Context, container types.Container) (domain.StorageNode, error) {
	c, err := n.c.ContainerInspect(ctx, container.ID)
	if err != nil {
		return nil, fmt.Errorf("inspect container id=%s: %w", c.ID, err)
	}

	accessKey, ok := getEnvValueByKey(c, EnvKeyMinioAccessKey)
	if !ok {
		return nil, fmt.Errorf("get minio access key failed")
	}
	secretKey, ok := getEnvValueByKey(c, EnvKeyMinioSecretKey)
	if !ok {
		return nil, fmt.Errorf("get minio secret key failed")
	}

	network, ok := c.NetworkSettings.Networks[n.networkIdentifier]
	if !ok {
		return nil, fmt.Errorf("cannot retrieve network=%s settings", n.networkIdentifier)
	}
	node, err := minio.NewNode(c.ID, net.JoinHostPort(network.IPAddress, n.nodePort), accessKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("new minio node: %w", err)
	}
	return node, nil
}

func getEnvValueByKey(container types.ContainerJSON, key string) (string, bool) {
	for _, env := range container.Config.Env {
		if strings.Contains(env, key+"=") {
			split := strings.Split(env, "=")
			return split[1], true
		}
	}
	return "", false
}
