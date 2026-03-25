package main

import (
	"context"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/reflect-homini/stora/internal/adapters/worker"
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/core/logger"
	"github.com/reflect-homini/stora/internal/core/otel"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger.Init("stora-worker")

	if err := config.Load(); err != nil {
		logger.Error(err)
		return 1
	}

	ctx := context.Background()
	otelShutdown, err := otel.InitSDK(ctx, config.Global.OTel)
	if err != nil {
		logger.Error(err)
		return 1
	}
	defer func() {
		if err := otelShutdown(ctx); err != nil {
			logger.Error(err)
		}
	}()

	w, err := worker.Setup()
	if err != nil {
		logger.Error(err)
		return 1
	}

	w.Run()
	return 0
}
