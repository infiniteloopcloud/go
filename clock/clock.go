package clock

import "time"

var _ Clock = Time{}
var _ Clock = Mock{}

// TODO add pkg lvl functions

var clock Clock = Time{}

type Clock interface {
	Now() time.Time
	Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) time.Time
	Unix(sec int64, nsec int64) time.Time
	UnixMilli(msec int64) time.Time
	UnixMicro(usec int64) time.Time
}

func SetClock(c Clock) {
	clock = c
}

func Now() time.Time {
	return clock.Now()
}
func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) time.Time {
	return clock.Date(year, month, day, hour, min, sec, nsec, loc)
}
func Unix(sec int64, nsec int64) time.Time {
	return clock.Unix(sec, nsec)
}
func UnixMilli(msec int64) time.Time {
	return clock.UnixMilli(msec)
}
func UnixMicro(usec int64) time.Time {
	return clock.UnixMicro(usec)
}

type Time struct{}

func (r Time) Now() time.Time {
	return time.Now().UTC()
}

func (r Time) Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, loc).UTC()
}

func (r Time) Unix(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec).UTC()
}

func (r Time) UnixMilli(msec int64) time.Time {
	return time.UnixMilli(msec).UTC()
}

func (r Time) UnixMicro(usec int64) time.Time {
	return time.UnixMicro(usec).UTC()
}

type Mock struct {
	Years  int
	Months int
	Days   int
}

func (m Mock) Now() time.Time {
	t := time.Now().AddDate(m.Years, m.Months, m.Days)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, t.Second(), t.Nanosecond(), time.UTC)
}

func (m Mock) Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, loc).UTC()
}

func (m Mock) Unix(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec).UTC()
}

func (m Mock) UnixMilli(msec int64) time.Time {
	return time.UnixMilli(msec).UTC()
}

func (m Mock) UnixMicro(usec int64) time.Time {
	return time.UnixMicro(usec).UTC()
}
