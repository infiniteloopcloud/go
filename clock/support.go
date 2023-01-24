package clock

import (
	"time"
)

func StartDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndDate(t time.Time) time.Time {
	return StartDate(t).AddDate(0, 0, 1).Add(-1 * time.Nanosecond)
}
