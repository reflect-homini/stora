package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/otel"
	"gorm.io/gorm"
)

type ProjectSummaryRepository interface {
	crud.Repository[ProjectSummary]
	FindLatest(ctx context.Context, projectID uuid.UUID) (ProjectSummary, error)
}

func NewProjectSummaryRepository(db *gorm.DB) *projectSummaryRepo {
	return &projectSummaryRepo{crud.NewRepository[ProjectSummary](db)}
}

type projectSummaryRepo struct {
	crud.Repository[ProjectSummary]
}

func (psr *projectSummaryRepo) FindLatest(ctx context.Context, projectID uuid.UUID) (ProjectSummary, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectSummaryRepository.FindLatest")
	defer span.End()

	db, err := psr.GetGormInstance(ctx)
	if err != nil {
		return ProjectSummary{}, err
	}

	var model ProjectSummary
	if err = db.
		Where("project_id", projectID).
		Order("generated_at DESC").
		First(&model).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model, nil
		}
		return ProjectSummary{}, ungerr.Wrap(err, "error querying latest project summary")
	}

	return model, nil
}
