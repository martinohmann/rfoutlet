package outlet

import (
	"sync"

	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// State describes the state of an outlet (on or off).
type State uint

const (
	// StateOff describes an outlet that is switched off.
	StateOff State = iota
	// StateOn describes an outlet that is switched on.
	StateOn
)

// Group is a group of outlets that can be switched as one.
type Group struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Outlets     []*Outlet `json:"outlets"`
}

// Outlet is an rf controlled outlet that can be switched on or off.
type Outlet struct {
	sync.Mutex
	ID          string             `json:"id"`
	DisplayName string             `json:"displayName"`
	CodeOn      uint64             `json:"-"`
	CodeOff     uint64             `json:"-"`
	Protocol    int                `json:"-"`
	PulseLength uint               `json:"-"`
	Schedule    *schedule.Schedule `json:"schedule"`
	State       State              `json:"state"`
}

// SetState sets the state of the outlet
func (o *Outlet) SetState(state State) {
	o.Lock()
	o.State = state
	o.Unlock()
}

// GetState returns the state of the outlet
func (o *Outlet) GetState() State {
	o.Lock()
	defer o.Unlock()
	return o.State
}

// getCodeForState returns the code to transmit to bring the outlet into state.
func (o *Outlet) getCodeForState(state State) uint64 {
	switch state {
	case StateOn:
		return o.CodeOn
	default:
		return o.CodeOff
	}
}
