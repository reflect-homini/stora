package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/ginkgo/pkg/server"
	"github.com/reflect-homini/stora/internal/appconstant"
	"github.com/reflect-homini/stora/internal/provider"
)

type Handlers struct {
	Auth           *AuthHandler
	Project        *ProjectHandler
	Entry          *EntryHandler
	ProjectSummary *ProjectSummaryHandler
}

func ProvideHandlers(services *provider.Services) *Handlers {
	return &Handlers{
		NewAuthHandler(services.Auth, services.OAuth, services.Session, services.User),
		&ProjectHandler{services.Project},
		&EntryHandler{services.Entry},
		&ProjectSummaryHandler{services.ProjectSummary},
	}
}

func getUserID(ctx *gin.Context) (uuid.UUID, error) {
	return server.GetFromContext[uuid.UUID](ctx, string(appconstant.ContextUserID))
}
