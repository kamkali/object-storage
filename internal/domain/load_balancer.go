package domain

import (
	"context"
)

type LoadBalancer interface {
	GetNode(ctx context.Context, key string) (StorageNode, error)
	Rebalance(ctx context.Context, nodes []StorageNode) error
}

//go:generate mockery --name=LoadBalancer
