package app

import (
	"context"
	"fmt"
	"github.com/PerfilievAlexandr/platform_common/pkg/closer"
	"github.com/PerfilievAlexandr/storage/internal/config"
	"log"
	"net/http"
	"sync"
)

type App struct {
	diProvider *diProvider
	httpServer *http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := a.runHttpServer(ctx)
		if err != nil {
			log.Fatal("failed to run HTTP server")
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initProvider,
		a.initHttpServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initProvider(_ context.Context) error {
	a.diProvider = newDiProvider()
	return nil
}

func (a *App) initHttpServer(ctx context.Context) error {
	a.diProvider.HttpHandler(ctx)
	a.httpServer = &http.Server{
		Addr:    a.diProvider.Config(ctx).HttpConfig.Address(),
		Handler: a.diProvider.httpHandler.InitRoutes(),
	}

	return nil
}

func (a *App) runHttpServer(ctx context.Context) error {
	log.Print(fmt.Sprintf("HTTP server is running on: %s", a.diProvider.Config(ctx).HttpConfig.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
