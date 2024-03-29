package manager

import (
	"fmt"

	"github.com/kamkalis/object-storage/internal/domain"
	"golang.org/x/net/context"
)

// StorageManager handles storage nodes refreshment — it discovers new nodes and handle them to the LoadBalancer
type StorageManager struct {
	domain.LoadBalancer
	domain.NodeDiscoverer
}

func NewStorageManager(lb domain.LoadBalancer, nodeDiscoverer domain.NodeDiscoverer) *StorageManager {
	return &StorageManager{LoadBalancer: lb, NodeDiscoverer: nodeDiscoverer}
}

func (s *StorageManager) RefreshNodes(ctx context.Context) error {
	nodes, err := s.DiscoverNodes(ctx)
	if err != nil {
		return fmt.Errorf("discover nodes: %w", err)
	}

	if err := s.Rebalance(ctx, nodes); err != nil {
		return fmt.Errorf("rebalance nodes: %w", err)
	}

	return nil
}
