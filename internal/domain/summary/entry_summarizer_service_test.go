package summary

import (
	"database/sql"
	"testing"
	"time"

	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/project"
	"github.com/stretchr/testify/assert"
)

func TestComputeTimeframeLabel(t *testing.T) {
	es := &entrySummarizer{}

	now := time.Now()

	createEntry := func(t time.Time) entry.Entry {
		e := entry.Entry{}
		e.CreatedAt = t
		return e
	}

	tests := []struct {
		name     string
		entries  []entry.Entry
		expected string
	}{
		{
			name: "Today",
			entries: []entry.Entry{
				createEntry(now),
				createEntry(now.Add(1 * time.Hour)),
			},
			expected: "Today",
		},
		{
			name: "Last 3 days",
			entries: []entry.Entry{
				createEntry(now.AddDate(0, 0, -2)),
				createEntry(now.AddDate(0, 0, -1)),
				createEntry(now),
			},
			expected: "Last 3 days",
		},
		{
			name: "January and March",
			entries: []entry.Entry{
				createEntry(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)),
				createEntry(time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)),
			},
			expected: "January and March",
		},
		{
			name: "January, February and March",
			entries: []entry.Entry{
				createEntry(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)),
				createEntry(time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)),
				createEntry(time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)),
			},
			expected: "January, February and March",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := es.computeTimeframeLabel(tt.entries)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestConstructPrompt(t *testing.T) {
	es := &entrySummarizer{}

	p := project.Project{Name: "Test Project"}
	entries := []entry.Entry{
		{Content: "Entry 1"},
	}
	entries[0].CreatedAt = time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	previousSummary := ProjectSummary{}
	previousSummary.SummaryMarkdown = sql.NullString{String: "Old Summary", Valid: true}

	prompt := es.constructPrompt(p, entries, "Today", previousSummary)

	assert.Contains(t, prompt.UserMessage, "Project: Test Project")
	assert.Contains(t, prompt.UserMessage, "Timeframe: Today")
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
