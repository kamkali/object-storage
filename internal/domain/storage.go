package domain

import (
	"errors"
	"io"

	"golang.org/x/net/context"
)

var (
	ErrObjNotFound = errors.New("object not found")
	ErrInvalidID   = errors.New("invalid ID")
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
	ID() string
	IsAlive(ctx context.Context) bool
	PutObject(ctx context.Context, o *Object) error
	GetObject(ctx context.Context, key string) (*Object, error)
}

//go:generate mockery --name=StorageNode

type Object struct {
	ID          string
	Content     io.Reader
	ContentType string
	Size        int64
}

type StorageService interface {
	PutObject(ctx context.Context, o *Object) error
	GetObject(ctx context.Context, key string) (*Object, error)
}

//go:generate mockery --name=StorageService
