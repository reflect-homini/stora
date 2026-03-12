package entry

import (
	"context"

	"github.com/itsLeonB/go-crud"
	"github.com/reflect-homini/stora/internal/core/otel"
)

type Service interface {
	Create(ctx context.Context, req NewEntryRequest) (EntryResponse, error)
}

func NewService(repo crud.Repository[Entry]) *service {
	return &service{repo}
}

type service struct {
	repo crud.Repository[Entry]
}

func (s *service) Create(ctx context.Context, req NewEntryRequest) (EntryResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "EntryService.Create")
	defer span.End()

	newEntry := Entry{
		ProjectID: req.ProjectID,
		Content:   req.Content,
	}

	insertedEntry, err := s.repo.Insert(ctx, newEntry)
	if err != nil {
		return EntryResponse{}, err
	}

	return EntryToResponse(insertedEntry), nil
}
