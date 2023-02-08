package service

import (
	"fmt"
	"regexp"

	"github.com/kamkalis/object-storage/internal/domain"
	"golang.org/x/net/context"
)

var idRe = regexp.MustCompile(`^[a-zA-Z0-9]{1,32}$`)

// StorageService handles and validates requests
type StorageService struct {
	manager domain.NodeManager
}

func NewStorage(manager domain.NodeManager) *StorageService {
	return &StorageService{
		manager: manager,
	}
}

func (s StorageService) PutObject(ctx context.Context, o *domain.Object) error {
	if !idRe.Match([]byte(o.ID)) {
		return domain.ErrInvalidID
	}

	node, err := s.manager.GetNode(ctx, o.ID)
	if err != nil {
		return fmt.Errorf("select node: %w", err)
	}

	if !node.IsAlive(ctx) {
		return fmt.Errorf("node %s offline", node.ID())
	}

	return node.PutObject(ctx, o)
}

func (s StorageService) GetObject(ctx context.Context, key string) (*domain.Object, error) {
	if !idRe.Match([]byte(key)) {
		return nil, domain.ErrInvalidID
	}

	node, err := s.manager.GetNode(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("select node: %w", err)
	}

	if !node.IsAlive(ctx) {
		return nil, fmt.Errorf("node %s offline", node.ID())
	}

	return node.GetObject(ctx, key)
}
