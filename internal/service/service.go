package service

import (
	"context"
	"github.com/PerfilievAlexandr/storage/internal/api/http/dtoHttpStorage"
)

type StorageService interface {
	Add(ctx context.Context, req dtoHttpStorage.AddRequest) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
