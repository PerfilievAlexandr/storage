package main

import (
	"context"
	"github.com/PerfilievAlexandr/storage/internal/app"
	"log"
)

func main() {
	ctx := context.Background()

	a, err := app.New(ctx)
	if err != nil {
		log.Fatal("failed to init app")
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatal("failed to run app")
	}
}
