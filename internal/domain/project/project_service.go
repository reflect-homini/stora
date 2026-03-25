package project

import (
	"context"
	"database/sql"
	"sort"

	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/otel"
	"github.com/reflect-homini/stora/internal/domain/entry"
)

type Service interface {
	// Public
	Create(ctx context.Context, req NewProjectRequest) (ProjectResponse, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]ProjectResponse, error)
	AddEntry(ctx context.Context, req entry.NewRequest) (entry.Response, error)
	GetEntriesAfter(ctx context.Context, userID, projectID, entryID uuid.UUID) ([]entry.Response, error)

	// Internal
	GetByID(ctx context.Context, id, userID uuid.UUID, forUpdate bool) (ProjectResponse, error)
}

type service struct {
	transactor crud.Transactor
	repo       crud.Repository[Project]
	entrySvc   entry.Service
}

func NewService(
	transactor crud.Transactor,
	repo crud.Repository[Project],
	entrySvc entry.Service,
) *service {
	return &service{
		transactor,
		repo,
		entrySvc,
	}
}

func (s *service) Create(ctx context.Context, req NewProjectRequest) (ProjectResponse, error) {
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

	return projectToResponse(insertedProject), nil
}

func (s *service) GetAll(ctx context.Context, userID uuid.UUID) ([]ProjectResponse, error) {
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

	return ezutil.MapSlice(projects, projectToResponse), nil
}

func (s *service) GetByID(ctx context.Context, id, userID uuid.UUID, forUpdate bool) (ProjectResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.GetByID")
	defer span.End()

	spec := crud.Specification[Project]{}
	spec.Model.ID = id
	spec.Model.UserID = userID
	spec.ForUpdate = forUpdate

	project, err := s.repo.FindFirst(ctx, spec)
	if err != nil {
		return ProjectResponse{}, err
	}

	if project.IsZero() {
		return ProjectResponse{}, ungerr.NotFoundError("project is not found")
	}

	return projectToResponse(project), nil
}

func (s *service) AddEntry(ctx context.Context, req entry.NewRequest) (entry.Response, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.AddEntry")
	defer span.End()

	var resp entry.Response
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		if _, err := s.GetByID(ctx, req.ProjectID, req.UserID, true); err != nil {
			return err
		}

		entry, err := s.entrySvc.Create(ctx, req)
		resp = entry
		return err
	})
	return resp, err
}

func (s *service) GetEntriesAfter(ctx context.Context, userID, projectID, entryID uuid.UUID) ([]entry.Response, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.GetEntriesAfter")
	defer span.End()

	if _, err := s.GetByID(ctx, projectID, userID, false); err != nil {
		return nil, err
	}

	return s.entrySvc.GetAfter(ctx, projectID, entryID)
}
