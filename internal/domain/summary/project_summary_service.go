package summary

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/logger"
	"github.com/reflect-homini/stora/internal/core/otel"
	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/project"
)

type ProjectSummaryService interface {
	GenerateDailySummaries(ctx context.Context) error
	GenerateDailySummary(ctx context.Context, projectID uuid.UUID) (ProjectSummary, error)
}

type projectSummaryService struct {
	repo            ProjectSummaryRepository
	projectRepo     crud.Repository[project.Project]
	entryRepo       entry.Repository
	entrySummarizer EntrySummarizerService
}

func NewProjectSummaryService(
	repo ProjectSummaryRepository,
	projectRepo crud.Repository[project.Project],
	entryRepo entry.Repository,
	entrySummarizer EntrySummarizerService,
) *projectSummaryService {
	return &projectSummaryService{
		repo,
		projectRepo,
		entryRepo,
		entrySummarizer,
	}
}

func (pss *projectSummaryService) GenerateDailySummary(ctx context.Context, projectID uuid.UUID) (ProjectSummary, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectSummaryService.GenerateDailySummary")
	defer span.End()

	spec := crud.Specification[project.Project]{}
	spec.Model.ID = projectID
	project, err := pss.projectRepo.FindFirst(ctx, spec)
	if err != nil {
		return ProjectSummary{}, err
	}
	if project.IsZero() {
		return ProjectSummary{}, ungerr.NotFoundError(fmt.Sprintf("project with ID %s not found", projectID))
	}

	summary, err := pss.generateSummary(ctx, project)
	if err != nil {
		return ProjectSummary{}, err
	}
	if summary.IsZero() {
		return ProjectSummary{}, nil
	}

	return pss.repo.Insert(ctx, summary)
}

func (pss *projectSummaryService) GenerateDailySummaries(ctx context.Context) error {
	ctx, span := otel.Tracer.Start(ctx, "ProjectSummaryService.GenerateDailySummary")
	defer span.End()

	projects, err := pss.projectRepo.FindAll(ctx, crud.Specification[project.Project]{})
	if err != nil {
		return err
	}

	if len(projects) < 1 {
		logger.Info("no projects found")
		return nil
	}

	newSummaries := make([]ProjectSummary, 0, len(projects))
	for _, project := range projects {
		summary, err := pss.generateSummary(ctx, project)
		if err != nil {
			span.RecordError(err)
			logger.Errorf("error generating summary for project ID %s: %v", project.ID, err)
			continue
		}
		if summary.ProjectID == uuid.Nil || !summary.SummaryText.Valid {
			continue
		}
		newSummaries = append(newSummaries, summary)
	}

	if len(newSummaries) < 1 {
		logger.Info("no new summaries to insert")
		return nil
	}

	_, err = pss.repo.InsertMany(ctx, newSummaries)
	return err
}

func (pss *projectSummaryService) generateSummary(ctx context.Context, project project.Project) (ProjectSummary, error) {
	latestSummary, entries, err := pss.getEntriesToSummarize(ctx, project.ID)
	if err != nil {
		return ProjectSummary{}, err
	}
	if len(entries) < 1 {
		logger.Infof("skipping summarization for project ID %s...", project.ID)
		return ProjectSummary{}, nil
	}

	return pss.entrySummarizer.Summarize(ctx, project, entries, latestSummary)
}

func (pss *projectSummaryService) getEntriesToSummarize(ctx context.Context, projectID uuid.UUID) (ProjectSummary, []entry.Entry, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectSummaryService.getEntriesToSummarize")
	defer span.End()

	latestSummary, err := pss.repo.GetLatest(ctx, projectID)
	if err != nil {
		return ProjectSummary{}, nil, err
	}

	if latestSummary.IsZero() {
		spec := crud.Specification[entry.Entry]{}
		spec.Model.ProjectID = projectID
		entries, err := pss.entryRepo.FindAll(ctx, spec)
		if err != nil {
			return ProjectSummary{}, nil, err
		}
		if len(entries) >= 5 {
			slices.Reverse(entries)
			return ProjectSummary{}, entries, nil
		}
		logger.Infof("insufficient entries for project ID %s, entries count: %d", projectID, len(entries))
		return ProjectSummary{}, nil, nil
	}

	newEntriesAfterLastSummary, err := pss.entryRepo.GetAfter(ctx, projectID, latestSummary.EndEntryID, 20)
	if err != nil {
		return ProjectSummary{}, nil, err
	}
	if len(newEntriesAfterLastSummary) >= 5 {
		return latestSummary, newEntriesAfterLastSummary, nil
	}

	logger.Infof("insufficient entries after last summary for project ID %s, entries count: %d", projectID, len(newEntriesAfterLastSummary))
	return ProjectSummary{}, nil, nil
}
