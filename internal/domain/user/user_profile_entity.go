package user

import (
	"database/sql"

	"github.com/itsLeonB/go-crud"
)

type UserProfile struct {
	crud.BaseEntity
	Name   string
	Avatar sql.NullString
}
