package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/ginkgo/pkg/server"
	"github.com/reflect-homini/stora/internal/appconstant"
	"github.com/reflect-homini/stora/internal/domain/project"
)

type EntryHandler struct {
	svc project.EntryService
}

func (eh *EntryHandler) HandleCreate() gin.HandlerFunc {
	return server.Handler("EntryHandler.HandleCreate", http.StatusCreated, func(ctx *gin.Context) (any, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		projectID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextProjectID))
		if err != nil {
			return nil, err
		}

		req, err := server.BindJSON[project.NewEntryRequest](ctx)
		if err != nil {
			return nil, err
		}

		req.UserID = userID
		req.ProjectID = projectID

		return eh.svc.Create(ctx.Request.Context(), req)
	})
}

func (eh *EntryHandler) HandleUpdate() gin.HandlerFunc {
	return server.Handler("EntryHandler.HandleUpdate", http.StatusOK, func(ctx *gin.Context) (any, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		projectID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextProjectID))
		if err != nil {
			return nil, err
		}

		entryID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextEntryID))
		if err != nil {
			return nil, err
		}

		req, err := server.BindJSON[project.UpdateEntryRequest](ctx)
		if err != nil {
			return nil, err
		}

		req.UserID = userID
		req.ProjectID = projectID
		req.ID = entryID

		return eh.svc.Update(ctx.Request.Context(), req)
	})
}

func (eh *EntryHandler) HandleDelete() gin.HandlerFunc {
	return server.Handler("EntryHandler.HandleDelete", http.StatusNoContent, func(ctx *gin.Context) (any, error) {
		userID, err := getUserID(ctx)
		if err != nil {
			return nil, err
		}

		projectID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextProjectID))
		if err != nil {
			return nil, err
		}

		entryID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextEntryID))
		if err != nil {
			return nil, err
		}

		req := project.DeleteEntryRequest{
			UserID:    userID,
			ProjectID: projectID,
			ID:        entryID,
		}

		return nil, eh.svc.Delete(ctx.Request.Context(), req)
	})
}
