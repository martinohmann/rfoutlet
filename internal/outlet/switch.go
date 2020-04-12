package outlet

import (
	"fmt"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

// Switcher defines the interface for an outlet switcher.
type Switcher interface {
	// Switch switches an outlet to the desired state.
	Switch(outlet *Outlet, state State) error
}

// Switch switches outlets by sending out codes using an gpio transmitter.
type Switch struct {
	Transmitter gpio.CodeTransmitter
}

// NewSwitch creates a new *Switch.
func NewSwitch(transmitter gpio.CodeTransmitter) *Switch {
	return &Switch{
		Transmitter: transmitter,
	}
}

// Switch switches an outlet to the provided state
func (s *Switch) Switch(o *Outlet, state State) error {
	if o.Protocol < 1 || o.Protocol > len(gpio.DefaultProtocols) {
		return fmt.Errorf("Protocol %d does not exist", o.Protocol)
	}

	proto := gpio.DefaultProtocols[o.Protocol-1]

	s.Transmitter.Transmit(o.getCodeForState(state), proto, o.PulseLength)
	o.SetState(state)

	return nil
}

// FakeSwitch can be used in tests,
type FakeSwitch struct {
	// Err is the error that should be returned by Switch. If non-nil, the
	// state of any passed in *Outlet will not be altered.
	Err error
}

// Switch implements Switcher.
//
// It will only set the outlet state or return the configured error.
func (s *FakeSwitch) Switch(outlet *Outlet, state State) error {
	if s.Err != nil {
		return s.Err
	}

	outlet.SetState(state)

	return nil
}
