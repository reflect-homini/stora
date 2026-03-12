package project

import (
	"github.com/itsLeonB/ezutil/v2"
	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/mapper"
)

func projectToResponse(p Project) ProjectResponse {
	return ProjectResponse{
		BaseDTO:     mapper.BaseToDTO(p.BaseEntity),
		UserID:      p.UserID,
		Name:        p.Name,
		Description: p.Description.String,

		// Relations
		Entries: ezutil.MapSlice(p.Entries, entry.EntryToResponse),
	}
}
