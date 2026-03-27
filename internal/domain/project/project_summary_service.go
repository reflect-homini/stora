package project

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/logger"
	"github.com/reflect-homini/stora/internal/core/otel"
)

type ProjectSummaryService interface {
	GetEntries(ctx context.Context, userID, projectID, summaryID uuid.UUID) ([]EntryResponse, error)
	GenerateAll(ctx context.Context) error
	Generate(ctx context.Context, projectID uuid.UUID) (ProjectSummary, error)
}

type projectSummaryService struct {
	repo            ProjectSummaryRepository
	entryRepo       EntryRepository
	entrySummarizer EntrySummarizerService
	projectSvc      ProjectService
}

func NewProjectSummaryService(
	repo ProjectSummaryRepository,
	entryRepo EntryRepository,
	entrySummarizer EntrySummarizerService,
	projectSvc ProjectService,
) *projectSummaryService {
	return &projectSummaryService{
		repo,
		entryRepo,
		entrySummarizer,
		projectSvc,
	}
}

func (s *projectSummaryService) GetEntries(ctx context.Context, userID, projectID, summaryID uuid.UUID) ([]EntryResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectSummaryService.GetEntries")
	defer span.End()

	if _, err := s.projectSvc.GetByID(ctx, projectID, userID, false); err != nil {
		return nil, err
	}

	spec := crud.Specification[ProjectSummary]{}
	spec.Model.ID = summaryID
	spec.Model.ProjectID = projectID
	summary, err := s.repo.FindFirst(ctx, spec)
	if err != nil {
		return nil, err
	}
	if summary.IsZero() {
		return nil, ungerr.NotFoundError(fmt.Sprintf("summary ID %s of project ID %s is not found", summaryID, projectID))
	}

	entries, err := s.entryRepo.GetBetween(ctx, projectID, summary.StartEntryID, summary.EndEntryID)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(entries, EntryToResponse), nil
}

func (pss *projectSummaryService) Generate(ctx context.Context, projectID uuid.UUID) (ProjectSummary, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectSummaryService.Generate")
	defer span.End()

	project, err := pss.projectSvc.Get(ctx, projectID)
	if err != nil {
		return ProjectSummary{}, err
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

func (pss *projectSummaryService) GenerateAll(ctx context.Context) error {
	ctx, span := otel.Tracer.Start(ctx, "ProjectSummaryService.GenerateAll")
	defer span.End()

	projects, err := pss.projectSvc.GetMany(ctx)
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

func (pss *projectSummaryService) generateSummary(ctx context.Context, project Project) (ProjectSummary, error) {
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

func (pss *projectSummaryService) getEntriesToSummarize(ctx context.Context, projectID uuid.UUID) (ProjectSummary, []Entry, error) {
	ctx, span := otel.Tracer.Start(ctx, "ProjectSummaryService.getEntriesToSummarize")
	defer span.End()

	latestSummary, err := pss.repo.FindLatest(ctx, projectID)
	if err != nil {
		return ProjectSummary{}, nil, err
	}

	if latestSummary.IsZero() {
		spec := crud.Specification[Entry]{}
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
