package timeframe

import (
	"time"
)

// ClassifyRelativeTimeframeSentence generates a relative timeframe sentence based on start and end times,
// entry count, and current time.
func ClassifyRelativeTimeframeSentence(start, end time.Time, entryCount int, now time.Time) string {
	start = start.UTC()
	end = end.UTC()
	now = now.UTC()

	if end.After(now) {
		return "Recently,"
	}

	delta := now.Sub(end)
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	endDay := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
	daySpan := max(int(endDay.Sub(startDay).Hours()/24)+1, 1)

	density := float64(entryCount) / float64(daySpan)
	isSparse := density < 0.5

	// Single-Day Cases
	if startDay.Equal(endDay) {
		nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		daysAgo := int(nowDay.Sub(endDay).Hours() / 24)

		if daysAgo == 0 {
			return "Today,"
		}
		if daysAgo == 1 {
			return "Yesterday,"
		}
	}

	// Multi-Day Recent Ranges
	if daySpan >= 2 && delta <= 7*24*time.Hour {
		if isSparse {
			return "On a few occasions over the past few days,"
		}
		return "Over the past few days,"
	}

	// Relative Time Buckets based on delta (now - end)
	const (
		day   = 24 * time.Hour
		week  = 7 * day
		month = 30 * day
	)

	if delta <= 1*day {
		// This covers cases where end was "yesterday" but maybe start was different.
		// Detailed rules say "Yesterday," (if not same-day)
		return "Yesterday,"
	}

	if delta <= 1*week {
		if isSparse {
			return "On a few occasions earlier this week,"
		}
		return "Earlier this week,"
	}

	if delta <= 2*week {
		if isSparse {
			return "On a few occasions over the past week,"
		}
		return "Last week,"
	}

	if delta <= 1*month {
		if isSparse {
			return "At a few points over the past month,"
		}
		return "A few weeks ago,"
	}

	if delta <= 3*month {
		if isSparse {
			return "At a few points over the past couple of months,"
		}
		return "A couple of months ago,"
	}

	// > 90 days
	if isSparse {
		return "At a few points a while ago,"
	}
	return "A while ago,"
}
