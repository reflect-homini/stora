package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/otel"
)

type EntryService interface {
	Create(ctx context.Context, req NewEntryRequest) (EntryResponse, error)
	Update(ctx context.Context, req UpdateEntryRequest) (EntryResponse, error)
	Delete(ctx context.Context, req DeleteEntryRequest) error
}

func NewEntryService(
	t crud.Transactor,
	repo EntryRepository,
	projectSvc ProjectService,
	projectSummaryRepo ProjectSummaryRepository,
) *entryService {
	return &entryService{
		t,
		repo,
		projectSvc,
		projectSummaryRepo,
	}
}

type entryService struct {
	transactor         crud.Transactor
	repo               EntryRepository
	projectSvc         ProjectService
	projectSummaryRepo ProjectSummaryRepository
}

func (s *entryService) Create(ctx context.Context, req NewEntryRequest) (EntryResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "EntryService.Create")
	defer span.End()

	var resp EntryResponse
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		if _, err := s.projectSvc.GetByID(ctx, req.ProjectID, req.UserID, true); err != nil {
			return err
		}

		newEntry := Entry{
			ProjectID: req.ProjectID,
			Content:   req.Content,
		}

		insertedEntry, err := s.repo.Insert(ctx, newEntry)
		if err != nil {
			return err
		}

		resp = EntryToResponse(insertedEntry)
		return nil
	})
	return resp, err
}

func (s *entryService) Update(ctx context.Context, req UpdateEntryRequest) (EntryResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "EntryService.Update")
	defer span.End()

	var resp EntryResponse
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.validateModification(ctx, req.UserID, req.ProjectID, req.ID); err != nil {
			return err
		}

		entry, err := s.getForUpdate(ctx, req.ProjectID, req.ID)
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

func (s *entryService) Delete(ctx context.Context, req DeleteEntryRequest) error {
	ctx, span := otel.Tracer.Start(ctx, "EntryService.Delete")
	defer span.End()

	return s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.validateModification(ctx, req.UserID, req.ProjectID, req.ID); err != nil {
			return err
		}

		entry, err := s.getForUpdate(ctx, req.ProjectID, req.ID)
		if err != nil {
			return err
		}

		return s.repo.Delete(ctx, entry)
	})
}

func (s *entryService) validateModification(ctx context.Context, userID, projectID, id uuid.UUID) error {
	summary, err := s.projectSummaryRepo.FindLatest(ctx, projectID)
	if err != nil {
		return err
	}

	if !summary.IsZero() && ezutil.CompareUUID(id, summary.EndEntryID) < 1 {
		return ungerr.ForbiddenError("cannot update summarized entry")
	}

	if _, err := s.projectSvc.GetByID(ctx, projectID, userID, true); err != nil {
		return err
	}

	return nil
}

func (s *entryService) getForUpdate(ctx context.Context, projectID, id uuid.UUID) (Entry, error) {
	spec := crud.Specification[Entry]{}
	spec.Model.ID = id
	spec.Model.ProjectID = projectID
	spec.ForUpdate = true
	entry, err := s.repo.FindFirst(ctx, spec)
	if err != nil {
		return Entry{}, err
	}
	if entry.IsZero() {
		return Entry{}, ungerr.NotFoundError(fmt.Sprintf("entry ID %s from project ID %s is not found", id, projectID))
	}

	return entry, nil
}
