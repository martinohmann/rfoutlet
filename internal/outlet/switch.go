package outlet

import (
	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

// Switcher defines the interface for an outlet switcher
type Switcher interface {
	Switch(*Outlet, State) error
}

// Switch type definition
type Switch struct {
	t gpio.CodeTransmitter
}

// NewSwitch creates a new switch
func NewSwitch(t gpio.CodeTransmitter) *Switch {
	return &Switch{t}
}

// Switch switches an outlet to the provided state
func (s *Switch) Switch(o *Outlet, state State) error {
	if err := s.t.Transmit(o.CodeForState(state), o.Protocol, o.PulseLength); err != nil {
		return err
	}

	o.SetState(state)

	return nil
}
