package user

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
)

type UserProfile struct {
	crud.BaseEntity
	UserID uuid.UUID
	Name   string
	Avatar sql.NullString
}
