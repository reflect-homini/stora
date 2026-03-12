package entry

import (
	"github.com/google/uuid"
	"github.com/reflect-homini/stora/internal/domain/dto"
)

type NewEntryRequest struct {
	UserID    uuid.UUID `json:"-" uri:"-"`
	ProjectID uuid.UUID `uri:"projectID" binding:"required"`
	Content   string    `json:"content" binding:"required,min=3"`
}

type EntryResponse struct {
	dto.BaseDTO
	ProjectID uuid.UUID `json:"projectId"`
	Content   string    `json:"content"`
}

type UpdateEntryRequest struct {
	UserID    uuid.UUID `json:"-" uri:"-"`
	ProjectID uuid.UUID `uri:"projectID" binding:"required"`
	ID        uuid.UUID `uri:"entryID" binding:"required"`
	Content   string    `json:"content" binding:"required,min=3"`
}

type DeleteEntryRequest struct {
	UserID    uuid.UUID `json:"-" uri:"-"`
	ProjectID uuid.UUID `uri:"projectID" binding:"required"`
	ID        uuid.UUID `uri:"entryID" binding:"required"`
}
