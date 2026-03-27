package entry

import "github.com/reflect-homini/stora/internal/domain/mapper"

func EntryToResponse(e Entry) Response {
	return Response{
		BaseDTO:   mapper.BaseToDTO(e.BaseEntity),
		ProjectID: e.ProjectID,
		Content:   e.Content,
	}
}
