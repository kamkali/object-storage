package domain

import (
	"context"
)

type LoadBalancer interface {
	GetNext(ctx context.Context) (StorageNode, error)
	ReBalance(ctx context.Context, nodes []StorageNode) error
}

//go:generate mockery --name=LoadBalancer
