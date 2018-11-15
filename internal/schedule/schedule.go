package schedule

import "time"

// Schedule type definition
type Schedule []Interval

// Enabled returns true if any of the intervals is enabled
func (s Schedule) Enabled() bool {
	for _, i := range s {
		if i.Enabled {
			return true
		}
	}

	return false
}

// Contains returns true if any of the intervals contains t
func (s Schedule) Contains(t time.Time) bool {
	for _, i := range s {
		if i.Contains(t) {
			return true
		}
	}

	return false
}
