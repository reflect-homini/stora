package entrymanip

import (
	"context"

	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/otel"
	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/project"
	"github.com/reflect-homini/stora/internal/domain/summary"
)

type Service interface {
	UpdateEntry(ctx context.Context, req entry.UpdateRequest) (entry.Response, error)
	DeleteEntry(ctx context.Context, req entry.DeleteRequest) error
}

type service struct {
	transactor         crud.Transactor
	projectSummaryRepo summary.ProjectSummaryRepository
	projectSvc         project.Service
	entrySvc           entry.Service
}

func NewService(
	transactor crud.Transactor,
	projectSummaryRepo summary.ProjectSummaryRepository,
	projectSvc project.Service,
	entrySvc entry.Service,
) *service {
	return &service{
		transactor,
		projectSummaryRepo,
		projectSvc,
		entrySvc,
	}
}

func (s *service) UpdateEntry(ctx context.Context, req entry.UpdateRequest) (entry.Response, error) {
	ctx, span := otel.Tracer.Start(ctx, "EntryManipulationService.UpdateEntry")
	defer span.End()

	var resp entry.Response
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		summary, err := s.projectSummaryRepo.GetLatest(ctx, req.ProjectID)
		if err != nil {
			return err
		}

		if ezutil.CompareUUID(req.ID, summary.EndEntryID) < 1 {
			return ungerr.ForbiddenError("cannot update summarized entry")
		}

		if _, err := s.projectSvc.GetByID(ctx, req.ProjectID, req.UserID, true); err != nil {
			return err
		}

		resp, err = s.entrySvc.Update(ctx, req)
		return err
	})
	return resp, err
}

func (s *service) DeleteEntry(ctx context.Context, req entry.DeleteRequest) error {
	ctx, span := otel.Tracer.Start(ctx, "EntryManipulationService.DeleteEntry")
	defer span.End()

	return s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		summary, err := s.projectSummaryRepo.GetLatest(ctx, req.ProjectID)
		if err != nil {
			return err
		}

		if ezutil.CompareUUID(req.ID, summary.EndEntryID) < 1 {
			return ungerr.ForbiddenError("cannot delete summarized entry")
		}

		if _, err := s.projectSvc.GetByID(ctx, req.ProjectID, req.UserID, true); err != nil {
			return err
		}

		return s.entrySvc.Delete(ctx, req)
	})
}
