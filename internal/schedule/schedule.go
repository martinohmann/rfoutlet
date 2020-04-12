package schedule

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Schedule type definition
type Schedule struct {
	sync.RWMutex
	intervals []Interval
}

// New create a new schedule
func New() *Schedule {
	return NewWithIntervals(make([]Interval, 0))
}

// NewWithIntervals create a new schedule seeded with intervals
func NewWithIntervals(intervals []Interval) *Schedule {
	return &Schedule{
		intervals: intervals,
	}
}

// Enabled returns true if any of the intervals is enabled
func (s *Schedule) Enabled() bool {
	if s == nil {
		return false
	}

	s.RLock()
	intervals := s.intervals
	s.RUnlock()

	for _, i := range intervals {
		if i.Enabled {
			return true
		}
	}

	return false
}

// Contains returns true if any of the intervals contains t
func (s *Schedule) Contains(t time.Time) bool {
	if s == nil {
		return false
	}

	s.RLock()
	intervals := s.intervals
	s.RUnlock()

	for _, i := range intervals {
		if i.Contains(t) {
			return true
		}
	}

	return false
}

// AddInterval adds an interval to the schedule of an outlet
func (s *Schedule) AddInterval(interval Interval) error {
	if interval.ID == "" {
		interval.ID = uuid.NewV4().String()
	}

	s.Lock()
	defer s.Unlock()

	for _, i := range s.intervals {
		if i.ID == interval.ID {
			return fmt.Errorf("interval %q already exists", interval.ID)
		}
	}

	s.intervals = append(s.intervals, interval)

	return nil
}

// UpdateInterval updates an interval of the schedule of an outlet. Will return
// an error if the interval does not exist.
func (s *Schedule) UpdateInterval(interval Interval) error {
	s.Lock()
	defer s.Unlock()

	for j, i := range s.intervals {
		if i.ID == interval.ID {
			s.intervals[j] = interval
			return nil
		}
	}

	return fmt.Errorf("interval %q does not exist", interval.ID)
}

// DeleteInterval deletes an interval of the schedule of an outlet. Will return
// an error if the interval does not exist.
func (s *Schedule) DeleteInterval(interval Interval) error {
	s.Lock()
	defer s.Unlock()

	for j, i := range s.intervals {
		if i.ID == interval.ID {
			s.intervals = append(s.intervals[:j], s.intervals[j+1:]...)
			return nil
		}
	}

	return fmt.Errorf("interval %q does not exist", interval.ID)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (s *Schedule) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &s.intervals)
}

// MarshalJSON implements the json.Marshaler interface
func (s *Schedule) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.intervals)
}
