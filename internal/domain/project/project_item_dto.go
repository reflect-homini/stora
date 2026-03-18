package project

import (
	"github.com/google/uuid"
	"github.com/reflect-homini/stora/internal/domain/dto"
)

type ItemType string

const (
	ItemTypeEntry   ItemType = "entry"
	ItemTypeSummary ItemType = "summary"
)

type ProjectItem struct {
	dto.BaseDTO
	ProjectID uuid.UUID `json:"projectId"`
	ItemType  ItemType  `json:"itemType"`
	Content   string    `json:"content"`

	// Summaries only
	EntriesCount int       `json:"entriesCount,omitzero"`
	EndEntryID   uuid.UUID `json:"endEntryId,omitzero"`
}
