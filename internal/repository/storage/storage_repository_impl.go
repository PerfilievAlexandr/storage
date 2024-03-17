package storageRepository

import (
	"context"
	"github.com/PerfilievAlexandr/storage/internal/api/http/dtoHttpStorage"
	"github.com/PerfilievAlexandr/storage/internal/errors"
	"github.com/PerfilievAlexandr/storage/internal/repository"
	"sync"
)

type repo struct {
	storage map[string]string
	mx      sync.RWMutex
}

func New(_ context.Context) repository.StorageRepository {
	return &repo{
		mx:      sync.RWMutex{},
		storage: make(map[string]string),
	}
}

func (r *repo) Add(_ context.Context, req dtoHttpStorage.AddRequest) error {
	r.mx.Lock()
	r.storage[req.Key] = req.Value
	r.mx.Unlock()

	return nil
}

func (r *repo) Get(_ context.Context, key string) (string, error) {
	r.mx.RLock()
	val, ok := r.storage[key]
	r.mx.RUnlock()

	if !ok {
		return "", customErrors.NotFound
	}

	return val, nil
}

func (r *repo) Delete(_ context.Context, key string) error {
	r.mx.Lock()
	delete(r.storage, key)
	r.mx.Unlock()

	return nil
}
