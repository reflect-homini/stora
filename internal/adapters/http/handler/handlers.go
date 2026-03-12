package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/itsLeonB/ginkgo/pkg/server"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/appconstant"
	"github.com/reflect-homini/stora/internal/core/logger"
	"github.com/reflect-homini/stora/internal/provider"
)

type Handlers struct {
	Auth    *AuthHandler
	Project *ProjectHandler
}

func ProvideHandlers(services *provider.Services) *Handlers {
	return &Handlers{
		NewAuthHandler(services.Auth, services.OAuth, services.Session, services.User),
		&ProjectHandler{services.Project},
	}
}

func getUserID(ctx *gin.Context) (uuid.UUID, error) {
	return server.GetFromContext[uuid.UUID](ctx, string(appconstant.ContextUserID))
}

func multiBind[T any](c *gin.Context, binders ...any) (T, error) {
	var obj T
	var err error
	for _, b := range binders {
		switch binder := b.(type) {
		case binding.BindingUri:
			err = c.ShouldBindUri(&obj)
		case binding.Binding:
			err = c.ShouldBindWith(&obj, binder)
		default:
			logger.Errorf("unsupported binding type: %T", binder)
		}
	}
	if err != nil {
		return obj, ungerr.Wrap(err, "failed to bind")
	}
	return obj, nil
}
