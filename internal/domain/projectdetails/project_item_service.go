package projectdetails

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/reflect-homini/stora/internal/core/otel"
	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/project"
	"github.com/reflect-homini/stora/internal/domain/summary"
)

type Service interface {
	GetByID(ctx context.Context, userID, projectID uuid.UUID) (project.ProjectResponse, error)
}

type service struct {
	projectSummaryRepo summary.ProjectSummaryRepository
	entryRepo          entry.Repository
	projectSvc         project.Service
}

func NewService(
	projectSummaryRepo summary.ProjectSummaryRepository,
	entryRepo entry.Repository,
	projectSvc project.Service,
) Service {
	return &service{
		projectSummaryRepo,
		entryRepo,
		projectSvc,
	}
}

func (s *service) GetByID(ctx context.Context, userID, projectID uuid.UUID) (project.ProjectResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectItemService.GetByID")
	defer span.End()

	proj, err := s.projectSvc.GetByID(ctx, projectID, userID, false)
	if err != nil {
		return project.ProjectResponse{}, err
	}
	items, err := s.getItems(ctx, projectID)
	if err != nil {
		return project.ProjectResponse{}, err
	}
	proj.Items = items
	return proj, nil
}

func (s *service) getItems(ctx context.Context, projectID uuid.UUID) ([]project.ProjectItem, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectItemService.getItems")
	defer span.End()

	spec := crud.Specification[summary.ProjectSummary]{}
	spec.Model.ProjectID = projectID
	summaries, err := s.projectSummaryRepo.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	if len(summaries) < 1 {
		return s.getEntries(ctx, projectID)
	}

	entries, err := s.entryRepo.GetAfter(ctx, projectID, summaries[0].EndEntryID, -1)
	if err != nil {
		return nil, err
	}

	items := make([]project.ProjectItem, 0, len(entries)+len(summaries))
	for _, entry := range entries {
		items = append(items, entryToItem(entry))
	}
	for _, summary := range summaries {
		items = append(items, summaryToItem(summary))
	}

	return items, nil
}

func (s *service) getEntries(ctx context.Context, projectID uuid.UUID) ([]project.ProjectItem, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectItemService.getEntries")
	defer span.End()

	spec := crud.Specification[entry.Entry]{}
	spec.Model.ProjectID = projectID
	entries, err := s.entryRepo.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(entries, entryToItem), nil
}
