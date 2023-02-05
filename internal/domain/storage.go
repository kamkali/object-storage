package domain

import (
    "errors"
    "github.com/google/uuid"
    "golang.org/x/net/context"
)

var (
    ErrObjNotFound = errors.New("object not found")
)

type StorageManager interface {
}

//go:generate mockery --name=StorageManager

type StorageServer interface {
}

//go:generate mockery --name=StorageServer

type Object struct {
    ID uuid.UUID
    // TODO: think about io.Writer
    Content []byte
}

type StorageService interface {
    PutObject(ctx context.Context, o *Object) error
    GetObject(ctx context.Context, id uuid.UUID) (*Object, error)
}

//go:generate mockery --name=StorageService
