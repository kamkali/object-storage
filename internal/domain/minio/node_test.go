package minio

import (
	"github.com/kamkalis/object-storage/internal/domain"
)

var _ domain.StorageNode = (*Node)(nil)

// TODO: integration test maybe
