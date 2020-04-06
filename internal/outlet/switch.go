package outlet

import (
	"fmt"

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
	if o.Protocol < 1 || o.Protocol > len(gpio.DefaultProtocols) {
		return fmt.Errorf("Protocol %d does not exist", o.Protocol)
	}

	proto := gpio.DefaultProtocols[o.Protocol-1]

	s.t.Transmit(o.CodeForState(state), proto, o.PulseLength)
	o.SetState(state)

	return nil
}
