package state

import (
	"sync"

	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// SwitchState describes the state of a switch.
type SwitchState uint

const (
	// SwitchStateOff describes a switch that is turned off.
	SwitchStateOff SwitchState = iota
	// SwitchStateOn describes a switch that is turned on.
	SwitchStateOn
)

// State describes the state of an outlet.
type State struct {
	sync.Mutex
	// SwitchState is the current state of the outlet switch.
	SwitchState SwitchState `json:"switchState"`
	// Schedule is the time switch schedule that is defined for the outlet.
	// Might be nil.
	Schedule *schedule.Schedule `json:"schedule"`
}

func (s *State) SetSwitchState(state SwitchState) {
	s.Lock()
	s.SwitchState = state
	s.Unlock()
}

func (s *State) GetSwitchState() SwitchState {
	s.Lock()
	defer s.Unlock()
	return s.SwitchState
}

func (s *State) SetSchedule(schedule *schedule.Schedule) {
	s.Lock()
	s.Schedule = schedule
	s.Unlock()
}

func (s *State) GetSchedule() *schedule.Schedule {
	s.Lock()
	defer s.Unlock()
	if s.Schedule == nil {
		s.Schedule = schedule.New()
	}
	return s.Schedule
}
