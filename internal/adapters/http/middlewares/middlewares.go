package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/ginkgo/pkg/middleware"
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/core/logger"
	"github.com/reflect-homini/stora/internal/domain/auth"
)

type Middlewares struct {
	Auth gin.HandlerFunc
	Err  gin.HandlerFunc
}

func Provide(configs config.App, authSvc auth.Service) *Middlewares {
	tokenCheckFunc := func(ctx *gin.Context, token string) (bool, map[string]any, error) {
		return authSvc.VerifyToken(ctx.Request.Context(), token)
	}

	middlewareProvider := middleware.NewMiddlewareProvider(logger.Global)
	authMiddleware := middlewareProvider.NewAuthMiddleware("Bearer", tokenCheckFunc)
	errorMiddleware := middlewareProvider.NewErrorMiddleware()

	return &Middlewares{
		authMiddleware,
		errorMiddleware,
	}
}
