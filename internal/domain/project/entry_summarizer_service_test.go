package project

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstructPrompt(t *testing.T) {
	es := &entrySummarizer{}

	p := Project{Name: "Test Project"}
	entries := []Entry{
		{Content: "Entry 1"},
	}
	entries[0].CreatedAt = time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	previousSummary := ProjectSummary{}
	previousSummary.SummaryMarkdown = sql.NullString{String: "Old Summary", Valid: true}

	prompt := es.constructPrompt(p, entries, previousSummary)

	assert.Contains(t, prompt.UserMessage, "Project: Test Project")
	assert.NotContains(t, prompt.UserMessage, "Timeframe:")
	assert.Contains(t, prompt.UserMessage, "Previous Summary:\nOld Summary")
	assert.Contains(t, prompt.UserMessage, "New Entries:")
	assert.Contains(t, prompt.UserMessage, "[2024-01-01 10:00:00] Entry 1")
}

func TestParseResponse(t *testing.T) {
	es := &entrySummarizer{}

	response := `{
  "summary_text": "Work was focused on implementing the new feature.",
  "summary_markdown": "## Key Themes\n- Feature development\n\n## Progress\n- Completed implementation",
  "insights_json": {
    "themes": ["Feature development"],
    "achievements": ["Completed implementation"],
    "challenges": [],
    "learnings": ["TDD helps catch edge cases early"],
    "skills": ["Go"],
    "impact": ["Faster delivery"]
  }
}`

	summaryText, summaryMarkdown, insightsJSON := es.parseResponse(response)

	assert.Equal(t, "Work was focused on implementing the new feature.", summaryText)
	assert.Equal(t, "## Key Themes\n- Feature development\n\n## Progress\n- Completed implementation", summaryMarkdown)
	assert.Contains(t, insightsJSON, "Feature development")
	assert.Contains(t, insightsJSON, "Completed implementation")
}

func TestParseResponseInvalidJson(t *testing.T) {
	es := &entrySummarizer{}

	response := "This is not valid JSON."
	summaryText, summaryMarkdown, insightsJSON := es.parseResponse(response)

	assert.Equal(t, "This is not valid JSON.", summaryText)
	assert.Equal(t, "", summaryMarkdown)
	assert.Equal(t, "{}", insightsJSON)
}
