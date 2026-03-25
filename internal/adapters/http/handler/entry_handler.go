package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/ginkgo/pkg/server"
	"github.com/reflect-homini/stora/internal/domain/appconstant"
	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/entrymanip"
)

type EntryHandler struct {
	entryManipSvc entrymanip.Service
}

func (eh *EntryHandler) HandleUpdateEntry() gin.HandlerFunc {
	return server.Handler("EntryHandler.HandleUpdateEntry", http.StatusOK, func(ctx *gin.Context) (any, error) {
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

		req, err := server.BindJSON[entry.UpdateRequest](ctx)
		if err != nil {
			return nil, err
		}

		req.UserID = userID
		req.ProjectID = projectID
		req.ID = entryID

		return eh.entryManipSvc.UpdateEntry(ctx, req)
	})
}

func (eh *EntryHandler) HandleDeleteEntry() gin.HandlerFunc {
	return server.Handler("EntryHandler.HandleDeleteEntry", http.StatusNoContent, func(ctx *gin.Context) (any, error) {
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

		req := entry.DeleteRequest{
			UserID:    userID,
			ProjectID: projectID,
			ID:        entryID,
		}

		return nil, eh.entryManipSvc.DeleteEntry(ctx, req)
	})
}
