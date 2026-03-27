package project

import "github.com/reflect-homini/stora/internal/domain/mapper"

func EntryToResponse(e Entry) EntryResponse {
	return EntryResponse{
		BaseDTO:   mapper.BaseToDTO(e.BaseEntity),
		ProjectID: e.ProjectID,
		Content:   e.Content,
	}
}
