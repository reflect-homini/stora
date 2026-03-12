package mapper

import (
	"github.com/itsLeonB/go-crud"
	"github.com/reflect-homini/stora/internal/domain/dto"
)

func BaseToDTO(be crud.BaseEntity) dto.BaseDTO {
	return dto.BaseDTO{
		ID:        be.ID,
		CreatedAt: be.CreatedAt,
		UpdatedAt: be.UpdatedAt,
	}
}
