package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/ginkgo/pkg/server"
	"github.com/reflect-homini/stora/internal/appconstant"
	"github.com/reflect-homini/stora/internal/core/logger"
	"github.com/reflect-homini/stora/internal/domain/project"
)

type ProjectSummaryHandler struct {
	svc project.ProjectSummaryService
}

func (ph *ProjectSummaryHandler) HandleGenerate() gin.HandlerFunc {
	return server.Handler("ProjectSummaryHandler.HandleGenerate", http.StatusOK, func(ctx *gin.Context) (any, error) {
		projectID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextProjectID))
		if err != nil {
			return nil, err
		}

		return ph.svc.Generate(ctx.Request.Context(), projectID)
	})
}

func (ph *ProjectSummaryHandler) HandleGenerateAll() gin.HandlerFunc {
	return server.Handler("ProjectSummaryHandler.HandleGenerateAll", http.StatusOK, func(ctx *gin.Context) (any, error) {
		go func() {
			if err := ph.svc.GenerateAll(context.Background()); err != nil {
				logger.Error(err)
			}
		}()
		return "job launched in background", nil
	})
}

func (ph *ProjectSummaryHandler) HandleGetEntries() gin.HandlerFunc {
	return server.Handler("ProjectSummaryHandler.HandleGetEntries", http.StatusOK, func(ctx *gin.Context) (any, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		projectID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextProjectID))
		if err != nil {
			return nil, err
		}

		summaryID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextSummaryID))
		if err != nil {
			return nil, err
		}

		return ph.svc.GetEntries(ctx.Request.Context(), userID, projectID, summaryID)
	})
}
