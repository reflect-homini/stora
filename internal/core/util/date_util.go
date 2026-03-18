package util

import (
	"time"

	"github.com/itsLeonB/ezutil/v2"
)

func GetStartAndEndOfToday() (time.Time, time.Time, error) {
	year, month, day := time.Now().Date()
	start, err := ezutil.GetStartOfDay(year, int(month), day)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := ezutil.GetEndOfDay(year, int(month), day)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return start, end, nil
}
