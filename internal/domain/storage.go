package domain

type StorageManager interface {
}

//go:generate mockery --name=StorageManager

type StorageServer interface {
}

//go:generate mockery --name=StorageServer

type StorageService interface {
}

//go:generate mockery --name=StorageService
