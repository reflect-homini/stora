package summary

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/reflect-homini/stora/internal/core/llm"
	"github.com/reflect-homini/stora/internal/core/otel"
	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/project"
)

type EntrySummarizerService interface {
	Summarize(ctx context.Context, project project.Project, entries []entry.Entry, previousSummary ProjectSummary) (ProjectSummary, error)
}

type entrySummarizer struct {
	llmSvc llm.LLMService
}

func NewEntrySummarizerService(
	llmSvc llm.LLMService,
) *entrySummarizer {
	return &entrySummarizer{
		llmSvc,
	}
}

func (es *entrySummarizer) Summarize(ctx context.Context, project project.Project, entries []entry.Entry, previousSummary ProjectSummary) (ProjectSummary, error) {
	ctx, span := otel.Tracer.Start(ctx, "EntrySummarizerService.Summarize")
	defer span.End()

	if len(entries) == 0 {
		return ProjectSummary{}, nil
	}

	timeframeLabel := es.computeTimeframeLabel(entries)
	periodStart := entries[0].CreatedAt
	periodEnd := entries[len(entries)-1].CreatedAt

	prompt := es.constructPrompt(project, entries, timeframeLabel, previousSummary)
	response, err := es.llmSvc.Prompt(ctx, prompt)
	if err != nil {
		return ProjectSummary{}, err
	}

	// Basic parsing assuming LLM returns markdown and JSON in a recognizable way
	// or we can just store the whole thing if the prompt is structured well.
	// For now, let's assume we want to split them if possible, or just store response in Markdown.
	// Requirements say: return summary_markdown and structured insights JSON.

	summaryMarkdown, insightsJSON := es.parseResponse(response)

	return ProjectSummary{
		ProjectID:       project.ID,
		SummaryMarkdown: sql.NullString{String: summaryMarkdown, Valid: true},
		InsightsJSON:    sql.NullString{String: insightsJSON, Valid: true},
		SummaryLevel:    DailyLevel,
		StartEntryID:    entries[0].ID,
		EndEntryID:      entries[len(entries)-1].ID,
		EntriesCount:    len(entries),
		TimeframeLabel:  timeframeLabel,
		PeriodStart:     periodStart,
		PeriodEnd:       periodEnd,
	}, nil
}

func (es *entrySummarizer) computeTimeframeLabel(entries []entry.Entry) string {
	if len(entries) == 0 {
		return ""
	}

	days := make(map[string]time.Time)
	months := make(map[time.Month]bool)
	var dates []time.Time

	for _, e := range entries {
		d := e.CreatedAt.Format("2006-01-02")
		if _, ok := days[d]; !ok {
			days[d] = e.CreatedAt
			dates = append(dates, e.CreatedAt)
		}
		months[e.CreatedAt.Month()] = true
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	if len(days) == 1 {
		return "Today"
	}

	// Check if consecutive
	isConsecutive := true
	if len(dates) > 1 {
		for i := 1; i < len(dates); i++ {
			diff := dates[i].Sub(dates[i-1])
			if diff > 24*time.Hour+1*time.Minute { // Allow some slack for DST/leap seconds if any, but simplified
				isConsecutive = false
				break
			}
		}
	}

	if isConsecutive {
		return fmt.Sprintf("Last %d days", len(days))
	}

	// Separate months
	var monthNames []string
	var sortedMonths []time.Month
	for m := range months {
		sortedMonths = append(sortedMonths, m)
	}
	sort.Slice(sortedMonths, func(i, j int) bool {
		return sortedMonths[i] < sortedMonths[j]
	})
	for _, m := range sortedMonths {
		monthNames = append(monthNames, m.String())
	}

	if len(monthNames) > 1 {
		return strings.Join(monthNames[:len(monthNames)-1], ", ") + " and " + monthNames[len(monthNames)-1]
	}

	return monthNames[0]
}

func (es *entrySummarizer) constructPrompt(project project.Project, entries []entry.Entry, timeframeLabel string, previousSummary ProjectSummary) llm.Prompt {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Project: %s\n", project.Name)
	if project.Description.Valid {
		fmt.Fprintf(&sb, "Project Description: %s\n", project.Description.String)
	}
	fmt.Fprintf(&sb, "Timeframe: %s\n", timeframeLabel)

	if previousSummary.SummaryMarkdown.Valid && previousSummary.SummaryMarkdown.String != "" {
		fmt.Fprintf(&sb, "\nPrevious Summary:\n%s\n", previousSummary.SummaryMarkdown.String)
	}

	fmt.Fprintf(&sb, "\nNew Entries:\n")
	for _, e := range entries {
		fmt.Fprintf(&sb, "[%s] %s\n", e.CreatedAt.Format(time.DateTime), e.Content)
	}

	return llm.Prompt{
		SystemMessage: es.getSystemPrompt(),
		UserMessage:   sb.String(),
	}
}

func (es *entrySummarizer) parseResponse(response string) (string, string) {
	// Simple parser: look for ```json ... ```
	jsonStart := strings.Index(response, "```json")
	if jsonStart == -1 {
		return response, "{}"
	}

	summaryMarkdown := strings.TrimSpace(response[:jsonStart])

	jsonEnd := strings.Index(response[jsonStart+7:], "```")
	if jsonEnd == -1 {
		return summaryMarkdown, "{}"
	}

	insightsJSON := response[jsonStart+7 : jsonStart+7+jsonEnd]
	insightsJSON = strings.TrimSpace(insightsJSON)

	return summaryMarkdown, insightsJSON
}

func (es *entrySummarizer) getSystemPrompt() string {
	return `
You are an assistant summarizing a user's work journal.

Your job is to summarize recent entries into a concise progress summary.

Use the provided timeframe label to describe when the work happened.
Do not describe inactive periods.

If entries occur in separate months, mention each month explicitly rather than describing continuous work.

Follow the output template exactly.

Context:

Project name
Project description (optional)
Timeframe label
Previous summary (optional)
New entries list with timestamps

Output format:

Part 1 — Markdown Summary

Sections must appear exactly as:

## Summary
Brief paragraph describing work completed during the timeframe.

## Key Themes
Bullet list of main topics.

## Progress
Bullet list of concrete accomplishments.

## Challenges
Bullet list of problems or blockers (optional).

## Learnings
Bullet list of insights gained.

Part 2 — Insights JSON

Valid JSON only.

Structure:

{
  "themes": [],
  "achievements": [],
  "challenges": [],
  "learnings": [],
  "skills": [],
  "impact": []
}

Rules:

- Do not invent technologies or concepts not present in entries.
- Extract high-level themes instead of specific tools when possible.
- Achievements must describe concrete work done.
- Skills should reflect demonstrated abilities.
`
}
