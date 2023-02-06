package minio

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/context"
	"time"
)

const (
	bucketName = "object-storage"
	// TODO: can be dynamic perhaps
)

type Node struct {
	id uuid.UUID
	c  *minio.Client
}

func NewNode(endpoint string, accessKeyID string, secretAccessKey string) (*Node, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("new minio node: %w", err)
	}

	n := &Node{
		id: uuid.New(),
		c:  minioClient,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := n.createBucket(ctx); err != nil {
		return nil, fmt.Errorf("create bucket: %w", err)
	}

	return n, nil
}

func (n *Node) createBucket(ctx context.Context) error {
	if err := n.c.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
		exists, errBucketExists := n.c.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (n *Node) ID() uuid.UUID {
	return n.id
}

func (n *Node) Addr(ctx context.Context) string {
	return n.c.EndpointURL().String()
}

func (n *Node) IsAlive(ctx context.Context) bool {
	return n.c.IsOnline()
}

func (n *Node) PutObject(ctx context.Context, o *domain.Object) error {
	_, err := n.c.PutObject(ctx, bucketName, o.ID.String(), o.Content, o.Size, minio.PutObjectOptions{ContentType: o.ContentType})
	if err != nil {
		return fmt.Errorf("put object=%s to minio: %w", o.ID.String(), err)
	}
	return nil
}

func (n *Node) GetObject(ctx context.Context, id uuid.UUID) (*domain.Object, error) {
	object, err := n.c.GetObject(ctx, bucketName, id.String(), minio.GetObjectOptions{})
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
		Size:    s.Size,
	}, nil
}
