package domain

import (
	"context"

	"github.com/google/uuid"
)

type LoadBalancer interface {
	GetNode(ctx context.Context, key uuid.UUID) (StorageNode, error)
	ReBalance(ctx context.Context, nodes []StorageNode) error
}

//go:generate mockery --name=LoadBalancer
