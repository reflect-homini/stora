package project

import (
	"github.com/reflect-homini/stora/internal/domain/mapper"
)

func projectToResponse(p Project) ProjectResponse {
	return ProjectResponse{
		BaseDTO:          mapper.BaseToDTO(p.BaseEntity),
		UserID:           p.UserID,
		Name:             p.Name,
		Description:      p.Description.String,
		LastInteractedAt: p.LastInteractedAt,
	}
}
