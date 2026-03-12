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
	Create(ctx context.Context, req NewProjectRequest) (ProjectResponse, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]ProjectResponse, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (ProjectResponse, error)
	AddEntry(ctx context.Context, req entry.NewEntryRequest) (entry.EntryResponse, error)
	UpdateEntry(ctx context.Context, req entry.UpdateEntryRequest) (entry.EntryResponse, error)
	DeleteEntry(ctx context.Context, req entry.DeleteEntryRequest) error
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

func (s *service) GetByID(ctx context.Context, id, userID uuid.UUID) (ProjectResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.GetByID")
	defer span.End()

	spec := crud.Specification[Project]{}
	spec.Model.ID = id
	spec.Model.UserID = userID
	spec.PreloadRelations = []string{"Entries"}
	project, err := s.getBySpec(ctx, spec)
	if err != nil {
		return ProjectResponse{}, err
	}

	return projectToResponse(project), nil
}

func (s *service) AddEntry(ctx context.Context, req entry.NewEntryRequest) (entry.EntryResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.AddEntry")
	defer span.End()

	var resp entry.EntryResponse
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		err := s.ensureUserOwned(ctx, req.UserID, req.ProjectID)
		if err != nil {
			return err
		}

		resp, err = s.entrySvc.Create(ctx, req)
		return err
	})
	return resp, err
}

func (s *service) UpdateEntry(ctx context.Context, req entry.UpdateEntryRequest) (entry.EntryResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.UpdateEntry")
	defer span.End()

	var resp entry.EntryResponse
	err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		err := s.ensureUserOwned(ctx, req.UserID, req.ProjectID)
		if err != nil {
			return err
		}

		resp, err = s.entrySvc.Update(ctx, req)
		return err
	})
	return resp, err
}

func (s *service) DeleteEntry(ctx context.Context, req entry.DeleteEntryRequest) error {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.DeleteEntry")
	defer span.End()

	return s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		err := s.ensureUserOwned(ctx, req.UserID, req.ProjectID)
		if err != nil {
			return err
		}

		return s.entrySvc.Delete(ctx, req)
	})
}

func (s *service) ensureUserOwned(ctx context.Context, userID, projectID uuid.UUID) error {
	spec := crud.Specification[Project]{}
	spec.Model.ID = projectID
	spec.Model.UserID = userID
	spec.ForUpdate = true
	_, err := s.getBySpec(ctx, spec)
	return err
}

func (s *service) getBySpec(ctx context.Context, spec crud.Specification[Project]) (Project, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.getBySpec")
	defer span.End()

	project, err := s.repo.FindFirst(ctx, spec)
	if err != nil {
		return Project{}, err
	}
	if project.IsZero() {
		return Project{}, ungerr.NotFoundError("project is not found")
	}
	return project, nil
}
