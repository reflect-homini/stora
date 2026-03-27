package projectdetails

import (
	"database/sql"
	"testing"
	"time"

	"github.com/reflect-homini/stora/internal/domain/summary"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeSummary(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Built feature X", "built feature X"},
		{"Google Cloud issue", "Google Cloud issue"},
		{"API integration", "API integration"},
		{"  .Fixed bug  ", "fixed bug"},
		{"Completed task.", "completed task."},
		{"Go development", "go development"},
		{"Already lowercase", "already lowercase"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalizeSummary(tt.input))
		})
	}
}

func TestSummaryToItem(t *testing.T) {
	now := time.Now()
	s := summary.ProjectSummary{
		SummaryText:  sql.NullString{String: "Fixed latency issues.", Valid: true},
		PeriodStart:  now.Add(-2 * time.Hour),
		PeriodEnd:    now.Add(-1 * time.Hour),
		EntriesCount: 5,
	}

	item := summaryToItem(s, now)

	// "Today," + " " + "fixed latency issues."
	assert.Equal(t, "Today, fixed latency issues.", item.Content)
}
