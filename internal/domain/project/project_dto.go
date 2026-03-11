package project

import (
	"github.com/google/uuid"
	"github.com/reflect-homini/stora/internal/domain/dto"
)

type NewProjectRequest struct {
	UserID      uuid.UUID `json:"-"`
	Name        string    `json:"name" binding:"required,min=3"`
	Description string    `json:"description"`
}

type ProjectResponse struct {
	dto.BaseDTO
	UserID      uuid.UUID `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
