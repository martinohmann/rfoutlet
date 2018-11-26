package outlet

import (
	"fmt"
	"sync"

	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// StateHandler defines the interface for a state handler, that loads and saves
// the state of an outlet.
type StateHandler interface {
	LoadState([]*Outlet) error
	SaveState([]*Outlet) error
}

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
	Name        string             `json:"name"`
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

// CodeForState returns the code to transmit to bring the outlet into state
func (o *Outlet) CodeForState(state State) uint64 {
	switch state {
	case StateOn:
		return o.CodeOn
	default:
		return o.CodeOff
	}
}

// Register registers an outlet to given name
func (m *Manager) Register(name string, outlet *Outlet) {
	m.Lock()
	m.outlets[name] = outlet
	m.Unlock()
}

// Get retrieves the outlet for given name
func (m *Manager) Get(name string) (*Outlet, error) {
	m.Lock()
	defer m.Unlock()

	o, ok := m.outlets[name]
	if !ok {
		return nil, fmt.Errorf("unknown outlet %q", name)
	}

	return o, nil
}

// Outlets returns a slice with all registered outlets
func (m *Manager) Outlets() []*Outlet {
	m.Lock()
	defer m.Unlock()

	s := make([]*Outlet, 0, len(m.outlets))
	for _, o := range m.outlets {
		s = append(s, o)
	}

	return s
}
