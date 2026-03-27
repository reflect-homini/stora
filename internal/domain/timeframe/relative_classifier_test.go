package timeframe

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClassifyRelativeTimeframeSentence(t *testing.T) {
	now := time.Date(2026, 3, 26, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		start      time.Time
		end        time.Time
		entryCount int
		expected   string
	}{
		{
			name:       "Today",
			start:      now.Add(-2 * time.Hour),
			end:        now.Add(-1 * time.Hour),
			entryCount: 1,
			expected:   "Today,",
		},
		{
			name:       "Yesterday",
			start:      now.Add(-25 * time.Hour),
			end:        now.Add(-24 * time.Hour),
			entryCount: 1,
			expected:   "Yesterday,",
		},
		{
			name:       "Earlier this week (3 days ago)",
			start:      now.Add(-73 * time.Hour),
			end:        now.Add(-72 * time.Hour),
			entryCount: 1,
			expected:   "Earlier this week,",
		},
		{
			name:       "Last week (10 days ago)",
			start:      now.Add(-241 * time.Hour),
			end:        now.Add(-240 * time.Hour),
			entryCount: 1,
			expected:   "Last week,",
		},
		{
			name:       "A few weeks ago (25 days ago)",
			start:      now.Add(-601 * time.Hour),
			end:        now.Add(-600 * time.Hour),
			entryCount: 1,
			expected:   "A few weeks ago,",
		},
		{
			name:       "A couple of months ago (60 days ago)",
			start:      now.Add(-1441 * time.Hour),
			end:        now.Add(-1440 * time.Hour),
			entryCount: 1,
			expected:   "A couple of months ago,",
		},
		{
			name:       "A while ago (120 days ago)",
			start:      now.Add(-2881 * time.Hour),
			end:        now.Add(-2880 * time.Hour),
			entryCount: 1,
			expected:   "A while ago,",
		},
		{
			name:       "Over the past few days (3-day span within 7 days)",
			start:      now.Add(-72 * time.Hour),
			end:        now.Add(-2 * time.Hour),
			entryCount: 5, // density > 0.5
			expected:   "Over the past few days,",
		},
		{
			name:       "Multi-day sparse",
			start:      now.Add(-72 * time.Hour),
			end:        now.Add(-2 * time.Hour),
			entryCount: 1, // density < 0.5 (1/4)
			expected:   "On a few occasions over the past few days,",
		},
		{
			name:       "Future Guard",
			start:      now.Add(1 * time.Hour),
			end:        now.Add(2 * time.Hour),
			entryCount: 1,
			expected:   "Recently,",
		},
		{
			name:       "Start == End but entryCount > 1",
			start:      now.Add(-2 * time.Hour),
			end:        now.Add(-2 * time.Hour),
			entryCount: 5,
			expected:   "Today,",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ClassifyRelativeTimeframeSentence(tt.start, tt.end, tt.entryCount, now)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
