package main

import (
	"context"
	"embed"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pressly/goose/v3"
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/core/logger"
	"github.com/reflect-homini/stora/internal/core/otel"
	"github.com/reflect-homini/stora/internal/provider"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	logger.Init("Job")

	if err := config.Load(); err != nil {
		logger.Error(err)
		exitCode = 1
		return
	}

	ctx := context.Background()
	otelShutdown, err := otel.InitSDK(ctx, config.Global.OTel)
	if err != nil {
		logger.Error("failed to initialize OTel SDK: %v", err)
		exitCode = 1
		return
	}
	defer func() {
		if err := otelShutdown(ctx); err != nil {
			logger.Errorf("error shutting down OTel SDK: %v", err)
		}
	}()

	dataSource, err := provider.ProvideDataSource()
	if err != nil {
		logger.Error(err)
		exitCode = 1
		return
	}

	goose.SetBaseFS(migrations)
	goose.SetLogger(logger.Global)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error(err)
		exitCode = 1
		return
	}

	if err := goose.Up(dataSource.SQL, "migrations"); err != nil {
		logger.Error(err)
		exitCode = 1
		return
	}
}
