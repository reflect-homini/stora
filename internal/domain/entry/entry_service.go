package entry

import (
	"context"

	"github.com/itsLeonB/go-crud"
	"github.com/reflect-homini/stora/internal/core/otel"
)

type Service interface {
	Create(ctx context.Context, req NewEntryRequest) (EntryResponse, error)
	Update(ctx context.Context, req UpdateEntryRequest) (EntryResponse, error)
	Delete(ctx context.Context, req DeleteEntryRequest) error
}

func NewService(
	transactor crud.Transactor,
	repo crud.Repository[Entry],
) *service {
	return &service{
		transactor,
		repo,
	}
}

type service struct {
	transactor crud.Transactor
	repo       crud.Repository[Entry]
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

func (s *service) Update(ctx context.Context, req UpdateEntryRequest) (EntryResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "EntryService.Update")
	defer span.End()

	var resp EntryResponse
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[Entry]{}
		spec.Model.ID = req.ID
		spec.Model.ProjectID = req.ProjectID
		entry, err := s.repo.FindFirst(ctx, spec)
		if err != nil {
			return err
		}

		entry.Content = req.Content

		entry, err = s.repo.Update(ctx, entry)
		if err != nil {
			return err
		}

		resp = EntryToResponse(entry)
		return nil
	})
	return resp, err
}

func (s *service) Delete(ctx context.Context, req DeleteEntryRequest) error {
	ctx, span := otel.Tracer.Start(ctx, "EntryService.Delete")
	defer span.End()

	return s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[Entry]{}
		spec.Model.ID = req.ID
		spec.Model.ProjectID = req.ProjectID
		entry, err := s.repo.FindFirst(ctx, spec)
		if err != nil {
			return err
		}

		return s.repo.Delete(ctx, entry)
	})
}
