package entry

import (
	"context"
	"fmt"

	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/otel"
)

type Service interface {
	Create(ctx context.Context, req NewRequest) (Response, error)
	Update(ctx context.Context, req UpdateRequest) (Response, error)
	Delete(ctx context.Context, req DeleteRequest) error
}

func NewService(
	t crud.Transactor,
	repo Repository,
) *service {
	return &service{
		t,
		repo,
	}
}

type service struct {
	transactor crud.Transactor
	repo       Repository
}

func (s *service) Create(ctx context.Context, req NewRequest) (Response, error) {
	ctx, span := otel.Tracer.Start(ctx, "EntryService.Create")
	defer span.End()

	newEntry := Entry{
		ProjectID: req.ProjectID,
		Content:   req.Content,
	}

	insertedEntry, err := s.repo.Insert(ctx, newEntry)
	if err != nil {
		return Response{}, err
	}

	return EntryToResponse(insertedEntry), nil
}

func (s *service) Update(ctx context.Context, req UpdateRequest) (Response, error) {
	ctx, span := otel.Tracer.Start(ctx, "EntryService.Update")
	defer span.End()

	var resp Response
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[Entry]{}
		spec.Model.ID = req.ID
		spec.Model.ProjectID = req.ProjectID
		spec.ForUpdate = true
		entry, err := s.repo.FindFirst(ctx, spec)
		if err != nil {
			return err
		}
		if entry.IsZero() {
			return ungerr.NotFoundError(fmt.Sprintf("entry ID %s from project ID %s is not found", req.ID, req.ProjectID))
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

func (s *service) Delete(ctx context.Context, req DeleteRequest) error {
	ctx, span := otel.Tracer.Start(ctx, "EntryService.Delete")
	defer span.End()

	return s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[Entry]{}
		spec.Model.ID = req.ID
		spec.Model.ProjectID = req.ProjectID
		spec.ForUpdate = true
		entry, err := s.repo.FindFirst(ctx, spec)
		if err != nil {
			return err
		}
		if entry.IsZero() {
			return ungerr.NotFoundError(fmt.Sprintf("entry ID %s from project ID %s is not found", req.ID, req.ProjectID))
		}

		return s.repo.Delete(ctx, entry)
	})
}
