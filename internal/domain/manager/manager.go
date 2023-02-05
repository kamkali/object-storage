package manager

import (
	"github.com/google/uuid"
	"github.com/kamkalis/object-storage/internal/domain"
	"golang.org/x/net/context"
)

type StorageManager struct {
	lb         domain.LoadBalancer
	discoverer domain.NodeDiscoverer
	nodes      []domain.StorageNode
}

func (s StorageManager) SelectNode(ctx context.Context, objectID uuid.UUID) (domain.StorageNode, error) {
	return nil, nil
}
