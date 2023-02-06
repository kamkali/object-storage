package domain

import (
	"errors"
	"io"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

var (
	ErrObjNotFound = errors.New("object not found")
)

type NodeDiscoverer interface {
	DiscoverNodes(ctx context.Context) ([]StorageNode, error)
}

//go:generate mockery --name=NodeDiscoverer

type NodeManager interface {
	NodeDiscoverer
	LoadBalancer
	RefreshNodes(ctx context.Context) error
}

//go:generate mockery --name=NodeManager

type StorageNode interface {
	ID() uuid.UUID
	Addr(ctx context.Context) string
	IsAlive(ctx context.Context) bool
	PutObject(ctx context.Context, o *Object) error
	GetObject(ctx context.Context, id uuid.UUID) (*Object, error)
}

//go:generate mockery --name=StorageNode

type Object struct {
	ID          uuid.UUID
	Content     io.Reader
	ContentType string
	Size        int64
}

type StorageService interface {
	PutObject(ctx context.Context, o *Object) error
	GetObject(ctx context.Context, id uuid.UUID) (*Object, error)
}

//go:generate mockery --name=StorageService
