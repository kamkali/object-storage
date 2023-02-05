package domain

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

var (
	ErrObjNotFound = errors.New("object not found")
)

type NodeManager interface {
	SelectNode(ctx context.Context, objectID uuid.UUID) (StorageNode, error)
}

//go:generate mockery --name=StorageManager

type StorageNode interface {
	Addr(ctx context.Context) string
	IsAlive(ctx context.Context) bool
	PutObject(ctx context.Context, o *Object) error
	GetObject(ctx context.Context, id uuid.UUID) (*Object, error)
}

//go:generate mockery --name=StorageNode

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

type NodeDiscoverer interface {
}

//go:generate mockery --name=NodeDiscoverer
