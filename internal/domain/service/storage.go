package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kamkalis/object-storage/internal/domain"
	"golang.org/x/net/context"
)

type StorageService struct {
	manager domain.NodeManager
}

func NewStorage(manager domain.NodeManager) *StorageService {
	return &StorageService{
		manager: manager,
	}
}

func (s StorageService) PutObject(ctx context.Context, o *domain.Object) error {
	node, err := s.manager.GetNode(ctx, o.ID)
	if err != nil {
		return fmt.Errorf("select node: %w", err)
	}

	if !node.IsAlive(ctx) {
		return fmt.Errorf("node %s offline", node.ID().String())
	}

	return node.PutObject(ctx, o)
}

func (s StorageService) GetObject(ctx context.Context, id uuid.UUID) (*domain.Object, error) {
	node, err := s.manager.GetNode(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("select node: %w", err)
	}

	if !node.IsAlive(ctx) {
		return nil, fmt.Errorf("node %s offline", node.ID().String())
	}

	return node.GetObject(ctx, id)
}
