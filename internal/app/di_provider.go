package app

import (
	"context"
	"github.com/PerfilievAlexandr/storage/internal/api/http"
	"github.com/PerfilievAlexandr/storage/internal/config"
	"github.com/PerfilievAlexandr/storage/internal/repository"
	storageRepository "github.com/PerfilievAlexandr/storage/internal/repository/storage"
	"github.com/PerfilievAlexandr/storage/internal/service"
	storageService "github.com/PerfilievAlexandr/storage/internal/service/storage"
	fileTransactionLoggerService "github.com/PerfilievAlexandr/storage/internal/service/transaction_logger"
	"log"
)

type diProvider struct {
	config                   *config.Config
	httpHandler              *http.Handler
	storageService           service.StorageService
	transactionLoggerService service.TransactionLoggerService
	storageRepository        repository.StorageRepository
}

func newDiProvider() *diProvider {
	return &diProvider{}
}

func (p *diProvider) Config(ctx context.Context) *config.Config {
	if p.config == nil {
		cfg, err := config.NewConfig(ctx)
		if err != nil {
			log.Fatal("failed to get config")
		}

		p.config = cfg
	}

	return p.config
}

func (p *diProvider) HttpHandler(ctx context.Context) *http.Handler {
	if p.httpHandler == nil {
		p.httpHandler = http.New(p.StorageService(ctx))
	}

	return p.httpHandler
}

func (p *diProvider) StorageRepository(ctx context.Context) repository.StorageRepository {
	if p.storageRepository == nil {
		p.storageRepository = storageRepository.New(ctx)
	}

	return p.storageRepository
}

func (p *diProvider) StorageService(ctx context.Context) service.StorageService {
	if p.storageService == nil {
		p.storageService = storageService.New(ctx, p.StorageRepository(ctx))
	}

	return p.storageService
}

func (p *diProvider) FileTransactionLoggerService(ctx context.Context) service.TransactionLoggerService {
	if p.transactionLoggerService == nil {
		transactionLoggerService, err := fileTransactionLoggerService.New(ctx, "transaction.log")
		if err != nil {
			log.Fatal(err)
		}

		p.transactionLoggerService = transactionLoggerService
	}

	return p.transactionLoggerService
}
