package summary

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"
	"github.com/reflect-homini/stora/internal/core/llm"
	"github.com/reflect-homini/stora/internal/core/logger"
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

	logger.Infof("summarizing project ID %s with %d entries...", project.ID, len(entries))

	startEntry := entries[0]
	endEntry := entries[len(entries)-1]

	prompt := es.constructPrompt(project, entries, previousSummary)
	response, err := es.llmSvc.Prompt(ctx, prompt)
	if err != nil {
		return ProjectSummary{}, err
	}

	summaryText, summaryMarkdown, insightsJSON := es.parseResponse(response)
	logger.Infof("finished summarizing project ID %s", project.ID)

	return ProjectSummary{
		ProjectID:       project.ID,
		SummaryText:     sql.NullString{String: summaryText, Valid: true},
		SummaryMarkdown: sql.NullString{String: summaryMarkdown, Valid: true},
		InsightsJSON:    sql.NullString{String: insightsJSON, Valid: true},
		SummaryLevel:    DailyLevel,
		StartEntryID:    startEntry.ID,
		EndEntryID:      endEntry.ID,
		EntriesCount:    len(entries),
		PeriodStart:     startEntry.CreatedAt,
		PeriodEnd:       endEntry.CreatedAt,
		GeneratedAt:     time.Now(),
	}, nil
}

func (es *entrySummarizer) constructPrompt(project project.Project, entries []entry.Entry, previousSummary ProjectSummary) llm.Prompt {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Project: %s\n", project.Name)
	if project.Description.Valid {
		fmt.Fprintf(&sb, "Project Description: %s\n", project.Description.String)
	}

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
				"description": "A brief prose paragraph summarising the work done. Must not contain any temporal or relative-time language.",
			},
			"summary_markdown": map[string]any{
				"type":        "string",
				"description": "Markdown containing Key Themes, Progress, Challenges and Learnings sections. Must not contain any temporal or relative-time language.",
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

Your job is to produce a structured JSON summary of the provided entries.

Do not describe inactive periods.

You must return a JSON object with exactly three fields:

1. "summary_text"
   A brief prose paragraph (2-4 sentences) describing the work done.
   Do not use markdown here.
   Do NOT reference time, dates, or recency (e.g. do not say "recently", "this week",
   "in the entries", "during this period", or any similar temporal language).

2. "summary_markdown"
   A markdown string containing exactly these sections (omit a section only if there
   is genuinely nothing to report):

   ## Key Themes
   Bullet list of main topics.

   ## Progress
   Bullet list of concrete accomplishments.

   ## Challenges
   Bullet list of problems or blockers.

   ## Learnings
   Bullet list of insights gained.

   Do NOT include a "## Summary" section — that content belongs in "summary_text".
   Do NOT reference time, dates, or recency anywhere in this field.

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
- Never use temporal or relative-time language anywhere in your output (e.g. "recently",
  "this week", "yesterday", "in recent entries", "during this period", "in March",
  "over the past", etc.). The timeframe is provided separately by the caller.
`
}
