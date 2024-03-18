package service

import (
	"context"
	"github.com/PerfilievAlexandr/storage/internal/api/http/dtoHttpStorage"
	"github.com/PerfilievAlexandr/storage/internal/domain"
)

type StorageService interface {
	Put(ctx context.Context, req dtoHttpStorage.AddRequest) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type TransactionLoggerService interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	ReadEvents() (chan<- domain.Event, chan<- error)
	Run()
}
