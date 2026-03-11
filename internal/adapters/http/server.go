package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/ezutil/v2/zerolog"
	"github.com/kroma-labs/sentinel-go/httpserver"
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/core/logger"
	"github.com/reflect-homini/stora/internal/provider"
)

func Setup(configs config.Config) (*httpserver.Server, func(), error) {
	providers, err := provider.All()
	if err != nil {
		return nil, nil, err
	}

	shutdownFunc := func() {
		if err := providers.Shutdown(); err != nil {
			logger.Error(err)
		}
	}

	gin.SetMode(configs.App.Env)
	r := gin.New()
	r.HandleMethodNotAllowed = true

	zerologger := zerolog.Instance(logger.Global)

	skipPaths := []string{"/ping", "/livez", "/readyz", "/metrics"}
	if err := setupSentinel(r, skipPaths, zerologger); err != nil {
		shutdownFunc()
		return nil, nil, err
	}

	RegisterRoutes(r, configs, providers.Services)

	httpCfg := httpserver.ProductionConfig()
	httpCfg.LoggerConfig = &httpserver.LoggerConfig{
		Logger:    zerologger,
		SkipPaths: skipPaths,
	}
	httpCfg.Addr = fmt.Sprintf(":%s", configs.App.Port)

	srv := httpserver.New(
		httpserver.WithConfig(httpCfg),
		httpserver.WithServiceName(configs.OTel.ServiceName),
		httpserver.WithHandler(r),
		httpserver.WithLogger(zerologger),
	)

	return srv, shutdownFunc, nil
}
