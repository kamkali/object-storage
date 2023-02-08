package minio

import (
	"fmt"

	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/context"
)

const (
	bucketName = "object-storage"
)

// Node implements StorageNode for a MinIO
type Node struct {
	id string
	c  *minio.Client
}

func NewNode(id, endpoint, accessKeyID, secretAccessKey string) (*Node, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("new minio node: %w", err)
	}

	n := &Node{
		id: id,
		c:  minioClient,
	}

	if err := n.createStorage(context.Background()); err != nil {
		return nil, fmt.Errorf("create bucket: %w", err)
	}

	return n, nil
}

func (n *Node) createStorage(ctx context.Context) error {
	if err := n.c.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
		exists, err := n.c.BucketExists(ctx, bucketName)
		if err != nil && !exists {
			return err
		}
	}
	return nil
}

func (n *Node) PutObject(ctx context.Context, o *domain.Object) error {
	_, err := n.c.PutObject(ctx, bucketName, o.ID, o.Content, o.Size, minio.PutObjectOptions{ContentType: o.ContentType})
	if err != nil {
		return fmt.Errorf("put object=%s to minio: %w", o.ID, err)
	}
	return nil
}

func (n *Node) GetObject(ctx context.Context, id string) (*domain.Object, error) {
	object, err := n.c.GetObject(ctx, bucketName, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object=%s from minio: %w", id, err)
	}

	s, err := object.Stat()
	if err != nil {
		e := minio.ToErrorResponse(err)
		if e.Code == "NoSuchKey" {
			return nil, domain.ErrObjNotFound
		}
		return nil, fmt.Errorf("get stats from minio object=%s: %w", id, err)
	}
	return &domain.Object{
		ID:          id,
		Content:     object,
		ContentType: s.ContentType,
		Size:        s.Size,
	}, nil
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) IsAlive(ctx context.Context) bool {
	return n.c.IsOnline()
}
