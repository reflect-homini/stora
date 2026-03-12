package project

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
	"github.com/reflect-homini/stora/internal/domain/entry"
)

type Project struct {
	crud.BaseEntity
	UserID      uuid.UUID
	Name        string
	Description sql.NullString

	// Relations
	Entries []entry.Entry
}
