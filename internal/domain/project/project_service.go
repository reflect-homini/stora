package project

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/otel"
)

type Service interface {
	Create(ctx context.Context, req NewProjectRequest) (ProjectResponse, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]ProjectResponse, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (ProjectResponse, error)
}

type service struct {
	repo crud.Repository[Project]
}

func NewService(repo crud.Repository[Project]) *service {
	return &service{repo}
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

	return ezutil.MapSlice(projects, projectToResponse), nil
}

func (s *service) GetByID(ctx context.Context, id, userID uuid.UUID) (ProjectResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectService.GetByID")
	defer span.End()

	spec := crud.Specification[Project]{}
	spec.Model.ID = id
	spec.Model.UserID = userID
	project, err := s.repo.FindFirst(ctx, spec)
	if err != nil {
		return ProjectResponse{}, err
	}
	if project.IsZero() {
		return ProjectResponse{}, ungerr.NotFoundError(fmt.Sprintf("project ID %s is not found", id))
	}

	return projectToResponse(project), nil
}
