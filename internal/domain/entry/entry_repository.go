package entry

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/otel"
	"gorm.io/gorm"
)

type Repository interface {
	crud.Repository[Entry]
	GetAfter(ctx context.Context, projectID, entryID uuid.UUID, limit int) ([]Entry, error)
	GetByDateRange(ctx context.Context, projectID uuid.UUID, start time.Time, end time.Time) ([]Entry, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{crud.NewRepository[Entry](db)}
}

type repository struct {
	crud.Repository[Entry]
}

func (r *repository) GetAfter(ctx context.Context, projectID, entryID uuid.UUID, limit int) ([]Entry, error) {
	ctx, span := otel.Tracer.Start(ctx, "entryRepository.GetAfter")
	defer span.End()

	db, err := r.GetGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	var models []Entry

	if err = db.Where("project_id", projectID).
		Where("id > ?", entryID).
		Order("created_at").
		Limit(limit).
		Find(&models).
		Error; err != nil {
		return nil, ungerr.Wrapf(err, "error querying entries after %s", entryID)
	}

	return models, nil
}

func (r *repository) GetByDateRange(ctx context.Context, projectID uuid.UUID, start time.Time, end time.Time) ([]Entry, error) {
	ctx, span := otel.Tracer.Start(ctx, "entryRepository.GetByDateRange")
	defer span.End()

	db, err := r.GetGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	var models []Entry

	if err = db.Where("project_id", projectID).
		Where("created_at >= ? AND created_at <= ?", start, end).
		Order("created_at").
		Find(&models).
		Error; err != nil {
		return nil, ungerr.Wrapf(err, "error querying entries between %s and %s", start, end)
	}

	return models, nil
}
