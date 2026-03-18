package entry

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/reflect-homini/stora/internal/core/otel"
)

type Service interface {
	Create(ctx context.Context, req NewEntryRequest) (EntryResponse, error)
	GetAfter(ctx context.Context, projectID, entryID uuid.UUID) ([]EntryResponse, error)
}

func NewService(repo Repository) *service {
	return &service{repo}
}

type service struct {
	repo Repository
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

func (s *service) GetAfter(ctx context.Context, projectID, entryID uuid.UUID) ([]EntryResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "EntryService.GetAfter")
	defer span.End()

	entries, err := s.repo.GetAfter(ctx, projectID, entryID, -1)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(entries, EntryToResponse), nil
}
