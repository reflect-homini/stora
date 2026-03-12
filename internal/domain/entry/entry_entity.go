package entry

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
)

type Entry struct {
	crud.BaseEntity
	ProjectID uuid.UUID
	Content   string
}
