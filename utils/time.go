package utils

import (
	"time"
)

var NowFunc = time.Now

// IsWeekend reports whether the given date falls on a weekend.
func IsWeekend(date time.Time) bool {
	day := date.Weekday()
	return day == time.Saturday || day == time.Sunday
}

func TimeAgo(duration time.Duration) time.Time {
	now := NowFunc()
	return now.Add(-1 * duration)
}

func Now() time.Time {
	return NowFunc()
}
