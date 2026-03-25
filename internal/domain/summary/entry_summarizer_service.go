package summary

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"
	"github.com/reflect-homini/stora/internal/core/llm"
	"github.com/reflect-homini/stora/internal/core/otel"
	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/project"
)

// llmSummaryResponse matches the JSON schema enforced via the LLM response format.
type llmSummaryResponse struct {
	SummaryText     string          `json:"summary_text"`
	SummaryMarkdown string          `json:"summary_markdown"`
	InsightsJSON    llmInsightsJSON `json:"insights_json"`
}

type llmInsightsJSON struct {
	Themes       []string `json:"themes"`
	Achievements []string `json:"achievements"`
	Challenges   []string `json:"challenges"`
	Learnings    []string `json:"learnings"`
	Skills       []string `json:"skills"`
	Impact       []string `json:"impact"`
}

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

	startEntry := entries[len(entries)-1]
	endEntry := entries[0]

	prompt := es.constructPrompt(project, entries, timeframeLabel, previousSummary)
	response, err := es.llmSvc.Prompt(ctx, prompt)
	if err != nil {
		return ProjectSummary{}, err
	}

	summaryText, summaryMarkdown, insightsJSON := es.parseResponse(response)

	return ProjectSummary{
		ProjectID:       project.ID,
		SummaryText:     sql.NullString{String: summaryText, Valid: true},
		SummaryMarkdown: sql.NullString{String: summaryMarkdown, Valid: true},
		InsightsJSON:    sql.NullString{String: insightsJSON, Valid: true},
		SummaryLevel:    DailyLevel,
		StartEntryID:    startEntry.ID,
		EndEntryID:      endEntry.ID,
		EntriesCount:    len(entries),
		TimeframeLabel:  timeframeLabel,
		PeriodStart:     startEntry.CreatedAt,
		PeriodEnd:       endEntry.CreatedAt,
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
			if diff > 24*time.Hour+1*time.Minute {
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
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:        "summary_response",
					Description: openai.String("Structured summary of work journal entries"),
					Schema:      es.responseJSONSchema(),
					Strict:      openai.Bool(true),
				},
			},
		},
	}
}

// responseJSONSchema returns the JSON Schema used to enforce the LLM output structure.
func (es *entrySummarizer) responseJSONSchema() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"summary_text": map[string]any{
				"type":        "string",
				"description": "A brief paragraph summarising the work done during the timeframe.",
			},
			"summary_markdown": map[string]any{
				"type":        "string",
				"description": "Markdown containing Key Themes, Progress, Challenges and Learnings sections.",
			},
			"insights_json": map[string]any{
				"type":        "object",
				"description": "Structured insights extracted from the entries.",
				"properties": map[string]any{
					"themes":       map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					"achievements": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					"challenges":   map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					"learnings":    map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					"skills":       map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
					"impact":       map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
				},
				"required":             []string{"themes", "achievements", "challenges", "learnings", "skills", "impact"},
				"additionalProperties": false,
			},
		},
		"required":             []string{"summary_text", "summary_markdown", "insights_json"},
		"additionalProperties": false,
	}
}

// parseResponse unmarshals the structured JSON response from the LLM.
// Returns (summaryText, summaryMarkdown, insightsJSON).
func (es *entrySummarizer) parseResponse(response string) (string, string, string) {
	var parsed llmSummaryResponse
	if err := json.Unmarshal([]byte(response), &parsed); err != nil {
		// Fallback: treat the whole response as summary text, leave the rest empty.
		return response, "", "{}"
	}

	insightsBytes, err := json.Marshal(parsed.InsightsJSON)
	if err != nil {
		insightsBytes = []byte("{}")
	}

	return parsed.SummaryText, parsed.SummaryMarkdown, string(insightsBytes)
}

func (es *entrySummarizer) getSystemPrompt() string {
	return `
You are an assistant summarizing a user's work journal.

Your job is to produce a structured JSON summary of recent entries.

Use the provided timeframe label to describe when the work happened.
Do not describe inactive periods.
If entries occur in separate months, mention each month explicitly rather than describing continuous work.

You must return a JSON object with exactly three fields:

1. "summary_text"
   A brief prose paragraph (2-4 sentences) describing the overall work done during the timeframe.
   Do not use markdown here.

2. "summary_markdown"
   A markdown string containing exactly these sections (omit a section only if there is genuinely nothing to report):

   ## Key Themes
   Bullet list of main topics.

   ## Progress
   Bullet list of concrete accomplishments.

   ## Challenges
   Bullet list of problems or blockers.

   ## Learnings
   Bullet list of insights gained.

   Do NOT include a "## Summary" section — that content belongs in "summary_text".

3. "insights_json"
   An object with these keys (each an array of strings):
   - themes
   - achievements
   - challenges
   - learnings
   - skills
   - impact

Rules:
- Do not invent technologies or concepts not present in the entries.
- Extract high-level themes instead of specific tools when possible.
- Achievements must describe concrete work done.
- Skills should reflect demonstrated abilities.
`
}
