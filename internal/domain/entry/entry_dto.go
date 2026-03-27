package entry

import (
	"github.com/google/uuid"
	"github.com/reflect-homini/stora/internal/domain/dto"
)

type NewRequest struct {
	UserID    uuid.UUID `json:"-"`
	ProjectID uuid.UUID `json:"-"`
	Content   string    `json:"content" binding:"required,min=3"`
}

type Response struct {
	dto.BaseDTO
	ProjectID uuid.UUID `json:"projectId"`
	Content   string    `json:"content"`
}

type UpdateRequest struct {
	UserID    uuid.UUID `json:"-"`
	ID        uuid.UUID `json:"-"`
	ProjectID uuid.UUID `json:"-"`
	Content   string    `json:"content" binding:"required,min=3"`
}

type DeleteRequest struct {
	UserID    uuid.UUID `json:"-"`
	ID        uuid.UUID `json:"-"`
	ProjectID uuid.UUID `json:"-"`
}
