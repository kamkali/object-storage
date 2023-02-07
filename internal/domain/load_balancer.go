package domain

import (
	"context"
)

type LoadBalancer interface {
	GetNode(ctx context.Context, key string) (StorageNode, error)
	ReBalance(ctx context.Context, nodes []StorageNode) error
}

//go:generate mockery --name=LoadBalancer
