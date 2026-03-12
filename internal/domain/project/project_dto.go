package project

import (
	"time"

	"github.com/google/uuid"
	"github.com/reflect-homini/stora/internal/domain/dto"
	"github.com/reflect-homini/stora/internal/domain/entry"
)

type NewProjectRequest struct {
	UserID      uuid.UUID `json:"-"`
	Name        string    `json:"name" binding:"required,min=3"`
	Description string    `json:"description"`
}

type ProjectResponse struct {
	dto.BaseDTO
	UserID           uuid.UUID `json:"userId"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	LastInteractedAt time.Time `json:"lastInteractedAt"`

	// Relations
	Entries []entry.EntryResponse `json:"entries"`
}
