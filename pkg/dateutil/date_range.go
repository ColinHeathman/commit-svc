package dateutil

import (
	"time"
)

// DateRange for optionally holding a start date and an end date
type DateRange struct {
	Start *time.Time
	End   *time.Time
}

// NewDateRange initializes a DateRange, and sets the end date to the end of the day for easier comparison
func NewDateRange(start *time.Time, end *time.Time) DateRange {
	if end != nil {
		endTime := EndTimeOfDate(*end)
		return DateRange{
			start,
			&endTime,
		}
	}

	return DateRange{
		start,
		end,
	}
}

func EndTimeOfDate(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 59, t.Location())
}
