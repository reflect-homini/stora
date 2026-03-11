package project

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
)

type Project struct {
	crud.BaseEntity
	UserID      uuid.UUID
	Name        string
	Description sql.NullString
}
