package service

import (
    "github.com/google/uuid"
    "github.com/kamkalis/object-storage/internal/domain"
    "golang.org/x/net/context"
)

type StorageService struct {
}

func NewStorage() *StorageService {
    return &StorageService{}
}

func (s StorageService) PutObject(ctx context.Context, o *domain.Object) error {

    return nil
}

func (s StorageService) GetObject(ctx context.Context, id uuid.UUID) (*domain.Object, error) {
    return nil, nil
}
