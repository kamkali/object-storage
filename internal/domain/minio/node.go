package minio

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/context"
)

type Node struct {
	id         uuid.UUID
	bucketName string

	c *minio.Client
}

func NewNode(endpoint string, accessKeyID string, secretAccessKey string, bucketName string) (*Node, error) {
	// TODO: should I create new client for each node?
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("new minio node: %w", err)
	}

	return &Node{id: uuid.New(), bucketName: bucketName, c: minioClient}, nil
}

func (m *Node) ID() uuid.UUID {
	return m.id
}

func (m *Node) Addr(ctx context.Context) string {
	return m.c.EndpointURL().String()
}

func (m *Node) IsAlive(ctx context.Context) bool {
	return m.c.IsOnline()
}

func (m *Node) PutObject(ctx context.Context, o *domain.Object) error {
	_, err := m.c.PutObject(ctx, m.bucketName, o.ID.String(), o.Content, int64(o.Size), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return fmt.Errorf("put object=%s to minio: %w", o.ID.String(), err)
	}
	return nil
}

func (m *Node) GetObject(ctx context.Context, id uuid.UUID) (*domain.Object, error) {
	object, err := m.c.GetObject(ctx, m.bucketName, id.String(), minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object=%s from minio: %w", id.String(), err)
	}

	s, err := object.Stat()
	if err != nil {
		return nil, fmt.Errorf("get stats from minio object=%s: %w", id.String(), err)
	}
	return &domain.Object{
		ID:      id,
		Content: object,
		Size:    int(s.Size),
	}, nil
}
