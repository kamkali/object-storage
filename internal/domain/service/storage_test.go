package service

import (
    "github.com/kamkalis/object-storage/internal/domain"
)

var _ domain.StorageService = (*StorageService)(nil)
