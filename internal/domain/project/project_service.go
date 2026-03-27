package project

import (
	"context"
	"database/sql"
	"slices"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/otel"
)

type ProjectService interface {
	// Public
	Create(ctx context.Context, req NewProjectRequest) (ProjectResponse, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]ProjectResponse, error)
	GetDetails(ctx context.Context, userID, projectID uuid.UUID) (ProjectResponse, error)

	// Internal
	GetByID(ctx context.Context, id, userID uuid.UUID, forUpdate bool) (ProjectResponse, error)
	Get(ctx context.Context, id uuid.UUID) (Project, error)
	GetMany(ctx context.Context) ([]Project, error)
}

type projectService struct {
	transactor         crud.Transactor
	repo               crud.Repository[Project]
	projectSummaryRepo ProjectSummaryRepository
	entryRepo          EntryRepository
}

func NewProjectService(
	transactor crud.Transactor,
	repo crud.Repository[Project],
	projectSummaryRepo ProjectSummaryRepository,
	entryRepo EntryRepository,
) *projectService {
	return &projectService{
		transactor,
		repo,
		projectSummaryRepo,
		entryRepo,
	}
}

func (s *projectService) Create(ctx context.Context, req NewProjectRequest) (ProjectResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.Create")
	defer span.End()

	newProject := Project{
		UserID: req.UserID,
		Name:   req.Name,
		Description: sql.NullString{
			String: req.Description,
			Valid:  req.Description != "",
		},
	}

	insertedProject, err := s.repo.Insert(ctx, newProject)
	if err != nil {
		return ProjectResponse{}, err
	}

	return projectToProjectResponse(insertedProject), nil
}

func (s *projectService) GetAll(ctx context.Context, userID uuid.UUID) ([]ProjectResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.GetAll")
	defer span.End()

	spec := crud.Specification[Project]{}
	spec.Model.UserID = userID
	projects, err := s.repo.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].LastInteractedAt.After(projects[j].LastInteractedAt)
	})

	return ezutil.MapSlice(projects, projectToProjectResponse), nil
}

func (s *projectService) GetDetails(ctx context.Context, userID, projectID uuid.UUID) (ProjectResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.GetDetails")
	defer span.End()

	proj, err := s.GetByID(ctx, projectID, userID, false)
	if err != nil {
		return ProjectResponse{}, err
	}
	items, err := s.getItems(ctx, projectID)
	if err != nil {
		return ProjectResponse{}, err
	}
	proj.Items = items
	return proj, nil
}

func (s *projectService) getItems(ctx context.Context, projectID uuid.UUID) ([]ProjectItem, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.getItems")
	defer span.End()

	spec := crud.Specification[ProjectSummary]{}
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

	// temp fix to reverse ordering
	slices.Reverse(summaries)

	now := time.Now()

	items := make([]ProjectItem, 0, len(entries)+len(summaries))
	for _, summary := range summaries {
		items = append(items, summaryToItem(summary, now))
	}
	for _, entry := range entries {
		items = append(items, entryToItem(entry))
	}

	return items, nil
}

func (s *projectService) getEntries(ctx context.Context, projectID uuid.UUID) ([]ProjectItem, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.getEntries")
	defer span.End()

	spec := crud.Specification[Entry]{}
	spec.Model.ProjectID = projectID
	entries, err := s.entryRepo.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	// temp fix to reverse ordering
	slices.Reverse(entries)

	return ezutil.MapSlice(entries, entryToItem), nil
}

func (s *projectService) GetByID(ctx context.Context, id, userID uuid.UUID, forUpdate bool) (ProjectResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.GetByID")
	defer span.End()

	spec := crud.Specification[Project]{}
	spec.Model.ID = id
	spec.Model.UserID = userID
	spec.ForUpdate = forUpdate
	project, err := s.getBySpec(ctx, spec)
	if err != nil {
		return ProjectResponse{}, err
	}

	return projectToProjectResponse(project), nil
}

func (s *projectService) Get(ctx context.Context, id uuid.UUID) (Project, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.Get")
	defer span.End()

	spec := crud.Specification[Project]{}
	spec.Model.ID = id
	project, err := s.getBySpec(ctx, spec)
	if err != nil {
		return Project{}, err
	}

	return project, nil
}

func (s *projectService) GetMany(ctx context.Context) ([]Project, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.GetMany")
	defer span.End()

	return s.repo.FindAll(ctx, crud.Specification[Project]{})
}

func (s *projectService) getBySpec(ctx context.Context, spec crud.Specification[Project]) (Project, error) {
	project, err := s.repo.FindFirst(ctx, spec)
	if err != nil {
		return Project{}, err
	}

	if project.IsZero() {
		return Project{}, ungerr.NotFoundError("project is not found")
	}

	return project, nil
}
