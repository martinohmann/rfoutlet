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
	var code uint64
	switch state {
	case StateOn:
		code = o.CodeOn
	case StateOff:
		code = o.CodeOff
	}

	if err := s.t.Transmit(code, o.Protocol, o.PulseLength); err != nil {
		return err
	}

	o.SetState(state)

	return nil
}
