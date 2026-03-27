package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/ginkgo/pkg/server"
	"github.com/reflect-homini/stora/internal/appconstant"
	"github.com/reflect-homini/stora/internal/domain/project"
)

type ProjectHandler struct {
	svc project.ProjectService
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

		return ph.svc.GetDetails(ctx.Request.Context(), userID, projectID)
	})
}
