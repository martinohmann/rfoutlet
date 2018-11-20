package context

import (
	"fmt"
	"sync"

	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/martinohmann/rfoutlet/internal/state"
	uuid "github.com/satori/go.uuid"
)

// Group type definition
type Group struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Outlets []*Outlet `json:"outlets"`
}

// Outlet type definition
type Outlet struct {
	sync.RWMutex
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	CodeOn      uint64            `json:"-"`
	CodeOff     uint64            `json:"-"`
	Protocol    int               `json:"-"`
	PulseLength uint              `json:"-"`
	Schedule    schedule.Schedule `json:"schedule"`
	State       state.SwitchState `json:"state"`
}

// GetSwitchState returns the switch state of an outlet (thread-safe)
func (o *Outlet) GetSwitchState() state.SwitchState {
	o.RLock()
	defer o.RUnlock()

	return o.State
}

// SetSwitchState sets the switch state of an outlet (thread-safe)
func (o *Outlet) SetSwitchState(state state.SwitchState) {
	o.Lock()
	defer o.Unlock()

	o.State = state
}

// GetSchedule returns the schedule of an outlet (thread-safe)
func (o *Outlet) GetSchedule() schedule.Schedule {
	o.RLock()
	defer o.RUnlock()

	return o.Schedule
}

// AddInterval adds an interval to the schedule of an outlet (thread-safe)
func (o *Outlet) AddInterval(interval schedule.Interval) error {
	if interval.ID == "" {
		interval.ID = uuid.NewV4().String()
	}

	o.Lock()
	defer o.Unlock()

	for _, i := range o.Schedule {
		if i.ID == interval.ID {
			return fmt.Errorf("interval with identifier %q already exists", interval.ID)
		}
	}

	o.Schedule = append(o.Schedule, interval)

	return nil
}

// UpdateInterval updates an interval of the schedule of an outlet
// (thread-safe). Will return an error if the interval does not exist.
func (o *Outlet) UpdateInterval(interval schedule.Interval) error {
	o.Lock()
	defer o.Unlock()

	for i, intv := range o.Schedule {
		if intv.ID != interval.ID {
			continue
		}

		o.Schedule[i] = interval

		return nil
	}

	return fmt.Errorf("interval with identifier %q does not exist", interval.ID)
}

// DeleteInterval deletes an interval of the schedule of an outlet
// (thread-safe). Will return an error if the interval does not exist.
func (o *Outlet) DeleteInterval(interval schedule.Interval) error {
	o.Lock()
	defer o.Unlock()

	for i, intv := range o.Schedule {
		if intv.ID == interval.ID {
			o.Schedule = append(o.Schedule[:i], o.Schedule[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("interval with identifier %q does not exist", interval.ID)
}
