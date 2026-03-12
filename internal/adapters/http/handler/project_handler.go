package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/ginkgo/pkg/server"
	"github.com/reflect-homini/stora/internal/domain/appconstant"
	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/project"
)

type ProjectHandler struct {
	svc project.Service
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

		return ph.svc.GetByID(ctx.Request.Context(), projectID, userID)
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

		req, err := server.BindJSON[entry.NewEntryRequest](ctx)
		if err != nil {
			return nil, err
		}

		req.UserID = userID
		req.ProjectID = projectID

		return ph.svc.AddEntry(ctx.Request.Context(), req)
	})
}
