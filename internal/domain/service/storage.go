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
	node, err := s.manager.SelectNode(ctx, o.ID)
	if err != nil {
		return fmt.Errorf("select node: %w", err)
	}

	if err := node.PutObject(ctx, o); err != nil {
		return fmt.Errorf("put object on the node: %w", err)
	}

	return nil
}

func (s StorageService) GetObject(ctx context.Context, id uuid.UUID) (*domain.Object, error) {
	node, err := s.manager.SelectNode(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("select node: %w", err)
	}

	object, err := node.GetObject(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get object on the node: %w", err)
	}

	return object, nil
}
