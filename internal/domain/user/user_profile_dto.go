package user

import "github.com/google/uuid"

type NewProfileRequest struct {
	UserID uuid.UUID
	Name   string
	Avatar string
}
