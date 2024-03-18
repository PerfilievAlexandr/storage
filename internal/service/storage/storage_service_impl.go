package storageService

import (
	"context"
	"github.com/PerfilievAlexandr/storage/internal/api/http/dtoHttpStorage"
	"github.com/PerfilievAlexandr/storage/internal/service"

	"github.com/PerfilievAlexandr/storage/internal/repository"
)

type storageService struct {
	storageRepository repository.StorageRepository
}

func New(
	_ context.Context,
	storageRepository repository.StorageRepository,
) service.StorageService {
	return &storageService{storageRepository}
}

func (s *storageService) Put(ctx context.Context, req dtoHttpStorage.AddRequest) error {
	err := s.storageRepository.Put(ctx, req)

	if err != nil {
		return err
	}

	return nil
}

func (s *storageService) Get(ctx context.Context, key string) (string, error) {
	val, err := s.storageRepository.Get(ctx, key)

	if err != nil {
		return "", err
	}

	return val, nil
}

func (s *storageService) Delete(ctx context.Context, key string) error {
	err := s.storageRepository.Delete(ctx, key)

	if err != nil {
		return err
	}

	return nil
}
