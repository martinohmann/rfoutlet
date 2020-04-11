package outlet

import (
	"sync"

	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// State defines an outlet switch state
type State uint

const (
	// StateOff defines the state for a disabled switch
	StateOff State = iota

	// SwitchStateOn defines the state for an enabled switch
	StateOn
)

// Outlet type definition
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

// Group type definition
type Group struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Outlets     []*Outlet `json:"outlets"`
}
