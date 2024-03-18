package storeProcessService

import (
	"context"
	"github.com/PerfilievAlexandr/storage/internal/api/http/dtoHttpStorage"
	"github.com/PerfilievAlexandr/storage/internal/domain"
	"github.com/PerfilievAlexandr/storage/internal/domain/enum"
	"github.com/PerfilievAlexandr/storage/internal/service"
)

type storeProcessService struct {
	storageService           service.StorageService
	transactionLoggerService service.TransactionLoggerService
}

func New(
	_ context.Context,
	storageService service.StorageService,
	transactionLoggerService service.TransactionLoggerService,
) service.StoreProcessService {
	return &storeProcessService{storageService, transactionLoggerService}
}

func (s *storeProcessService) Put(ctx context.Context, req dtoHttpStorage.AddRequest) error {
	err := s.storageService.Put(ctx, req)
	s.transactionLoggerService.WritePut(req.Key, req.Value)

	if err != nil {
		return err
	}

	return nil
}

func (s *storeProcessService) Get(ctx context.Context, key string) (string, error) {
	val, err := s.storageService.Get(ctx, key)

	if err != nil {
		return "", err
	}

	return val, nil
}

func (s *storeProcessService) Delete(ctx context.Context, key string) error {
	err := s.storageService.Delete(ctx, key)
	s.transactionLoggerService.WriteDelete(key)

	if err != nil {
		return err
	}

	return nil
}

func (s *storeProcessService) SynchronizeStorage(ctx context.Context) {
	logger := s.transactionLoggerService

	events, errors := logger.ReadEvents()

	var event domain.Event
	var ok = true
	var err error

	for err == nil && ok {
		select {
		case err, ok = <-errors:
		case event, ok = <-events:
			switch event.EventType {
			case enum.EventDelete:
				err = s.storageService.Delete(ctx, event.Key)
			case enum.EventPut:
				err = s.storageService.Put(ctx, dtoHttpStorage.AddRequest{Key: event.Key, Value: event.Value})
			}
		}
	}

	logger.Run()
}
