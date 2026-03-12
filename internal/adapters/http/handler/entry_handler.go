package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/ginkgo/pkg/server"
	"github.com/reflect-homini/stora/internal/domain/appconstant"
	"github.com/reflect-homini/stora/internal/domain/entry"
)

type EntryHandler struct {
	svc entry.Service
}

func (eh *EntryHandler) HandleCreate() gin.HandlerFunc {
	return server.Handler("EntryHandler.HandleCreate", http.StatusCreated, func(ctx *gin.Context) (any, error) {
		projectID, err := server.GetRequiredPathParam[uuid.UUID](ctx, string(appconstant.ContextProjectID))
		if err != nil {
			return nil, err
		}

		req, err := server.BindJSON[entry.NewEntryRequest](ctx)
		if err != nil {
			return nil, err
		}

		req.ProjectID = projectID

		return eh.svc.Create(ctx, req)
	})
}
