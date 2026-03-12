package user

import (
	"github.com/google/uuid"
	"github.com/reflect-homini/stora/internal/domain/dto"
)

type NewProfileRequest struct {
	UserID uuid.UUID
	Name   string
	Avatar string
}

type ProfileResponse struct {
	dto.BaseDTO
	UserID uuid.UUID `json:"userId"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
}
