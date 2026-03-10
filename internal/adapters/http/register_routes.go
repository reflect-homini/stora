package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kroma-labs/sentinel-go/httpserver"
	sentinelGin "github.com/kroma-labs/sentinel-go/httpserver/adapters/gin"
	"github.com/reflect-homini/stora/internal/adapters/http/handler"
	"github.com/reflect-homini/stora/internal/adapters/http/middlewares"
	"github.com/reflect-homini/stora/internal/adapters/http/routes"
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/provider"
)

func RegisterRoutes(router *gin.Engine, configs config.Config, services *provider.Services) {
	handlers := handler.ProvideHandlers(services)
	middlewares := middlewares.Provide(configs.App, services.Auth)

	router.Use(middlewares.Err)

	sentinelGin.RegisterHealth(router, httpserver.NewHealthHandler())

	if configs.App.Env != "release" {
		sentinelGin.RegisterPprof(router, httpserver.DefaultPprofConfig())
	}

	routes.RegisterBaseRoutes(router)
	routes.RegisterAPIRoutes(router, handlers, middlewares.Auth)
}
