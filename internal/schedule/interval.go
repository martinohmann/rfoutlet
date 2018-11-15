package schedule

import "time"

// Interval type definition
type Interval struct {
	Enabled  bool
	Weekdays []time.Weekday
	From     DayTime
	To       DayTime
}

// Contains returns true if interval is enabled and t lies within
func (i Interval) Contains(t time.Time) bool {
	if !i.Enabled || !i.enabledOn(t.Weekday()) {
		return false
	}

	dt := NewDayTime(t.Hour(), t.Minute())

	if i.From.After(i.To) {
		return !dt.Between(i.To, i.From)
	}

	return dt.BetweenInclusive(i.From, i.To)
}

func (i Interval) enabledOn(weekday time.Weekday) bool {
	for _, wd := range i.Weekdays {
		if wd == weekday {
			return true
		}
	}

	return false
}
