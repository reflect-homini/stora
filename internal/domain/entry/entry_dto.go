package entry

import (
	"github.com/google/uuid"
	"github.com/reflect-homini/stora/internal/domain/dto"
)

type NewEntryRequest struct {
	UserID    uuid.UUID `json:"-"`
	ProjectID uuid.UUID `json:"-"`
	Content   string    `json:"content" binding:"required,min=3"`
}

type EntryResponse struct {
	dto.BaseDTO
	ProjectID uuid.UUID `json:"projectId"`
	Content   string    `json:"content"`
}
