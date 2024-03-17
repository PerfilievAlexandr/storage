package config

import (
	"context"
	httpConfig "github.com/PerfilievAlexandr/storage/internal/config/http"
	config "github.com/PerfilievAlexandr/storage/internal/config/interface"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	HttpConfig config.HttpServerConfig
}

func NewConfig(_ context.Context) (*Config, error) {
	httpCfg, err := httpConfig.NewHttpConfig()
	if err != nil {
		log.Fatal("failed to config")
	}

	return &Config{
		HttpConfig: httpCfg,
	}, nil
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
