package appconstant

import "time"

// Relative Time Buckets based on delta (now - end)
const (
	Day   = 24 * time.Hour
	Week  = 7 * Day
	Month = 30 * Day
)
