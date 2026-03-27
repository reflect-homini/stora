package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/ginkgo/pkg/server"
	"github.com/reflect-homini/stora/internal/appconstant"
	"github.com/reflect-homini/stora/internal/core/logger"
	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/project"
	"github.com/reflect-homini/stora/internal/domain/projectdetails"
	"github.com/reflect-homini/stora/internal/domain/summary"
)

type ProjectHandler struct {
	svc        project.Service
	summarySvc summary.ProjectSummaryService
	detailsSvc projectdetails.Service
	entrySvc   entry.Service
}

func (ph *ProjectHandler) HandleCreate() gin.HandlerFunc {
	return server.Handler("ProjectHandler.HandleCreate", http.StatusCreated, func(ctx *gin.Context) (any, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		req, err := server.BindJSON[project.NewProjectRequest](ctx)
		if err != nil {
			return nil, err
		}

		req.UserID = userID

		return ph.svc.Create(ctx.Request.Context(), req)
	})
}

func (ph *ProjectHandler) HandleGetAll() gin.HandlerFunc {
	return server.Handler("ProjectHandler.HandleGetAll", http.StatusOK, func(ctx *gin.Context) (any, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		return ph.svc.GetAll(ctx.Request.Context(), userID)
	})
}

func (ph *ProjectHandler) HandleGetByID() gin.HandlerFunc {
	return server.Handler("ProjectHandler.HandleGetByID", http.StatusOK, func(ctx *gin.Context) (any, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		projectID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextProjectID))
		if err != nil {
			return nil, err
		}

		return ph.detailsSvc.GetByID(ctx.Request.Context(), userID, projectID)
	})
}

func (ph *ProjectHandler) HandleAddEntry() gin.HandlerFunc {
	return server.Handler("ProjectHandler.HandleAddEntry", http.StatusCreated, func(ctx *gin.Context) (any, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		projectID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextProjectID))
		if err != nil {
			return nil, err
		}

		req, err := server.BindJSON[entry.NewRequest](ctx)
		if err != nil {
			return nil, err
		}

		req.UserID = userID
		req.ProjectID = projectID

		return ph.svc.AddEntry(ctx.Request.Context(), req)
	})
}

func (ph *ProjectHandler) HandleGenerateSummary() gin.HandlerFunc {
	return server.Handler("ProjectHandler.HandleGenerateSummary", http.StatusOK, func(ctx *gin.Context) (any, error) {
		projectID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextProjectID))
		if err != nil {
			return nil, err
		}

		return ph.summarySvc.GenerateDailySummary(ctx.Request.Context(), projectID)
	})
}

func (ph *ProjectHandler) HandleGenerateSummaries() gin.HandlerFunc {
	return server.Handler("ProjectHandler.HandleGenerateSummaries", http.StatusOK, func(ctx *gin.Context) (any, error) {
		go func() {
			if err := ph.summarySvc.GenerateDailySummaries(context.Background()); err != nil {
				logger.Error(err)
			}
		}()
		return "job launched in background", nil
	})
}

func (ph *ProjectHandler) HandleGetSummaryEntries() gin.HandlerFunc {
	return server.Handler("ProjectHandler.HandleGetSummaryEntries", http.StatusOK, func(ctx *gin.Context) (any, error) {
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

		return ph.summarySvc.GetEntries(ctx.Request.Context(), userID, projectID, summaryID)
	})
}
