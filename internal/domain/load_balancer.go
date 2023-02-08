package domain

import (
	"context"
)

type LoadBalancer interface {
	// GetNode selects a server for a given key
	GetNode(ctx context.Context, key string) (StorageNode, error)
	// Rebalance performs a redistribution of the servers
	Rebalance(ctx context.Context, nodes []StorageNode) error
}

//go:generate mockery --name=LoadBalancer
