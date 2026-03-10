package main

import (
	"context"
	"os"

	"github.com/reflect-homini/stora/internal/adapters/http"
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/core/logger"
	"github.com/reflect-homini/stora/internal/core/otel"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	logger.Init("dapur")

	if err := config.Load(); err != nil {
		logger.Error(err)
		exitCode = 1
		return
	}

	ctx := context.Background()
	otelShutdown, err := otel.InitSDK(ctx, config.Global.OTel)
	if err != nil {
		logger.Error(err)
		exitCode = 1
		return
	}
	defer func() {
		if err := otelShutdown(ctx); err != nil {
			logger.Error(err)
		}
	}()

	srv, _, err := http.Setup(*config.Global)
	if err != nil {
		logger.Error(err)
		exitCode = 1
		return
	}
	// defer shutdownFunc()

	if err := srv.ListenAndServe(ctx); err != nil {
		logger.Error(err)
		exitCode = 1
	}
}
