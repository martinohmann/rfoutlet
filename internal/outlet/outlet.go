package outlet

import (
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/gpio"
)

const (
	StateUnknown = iota
	StateOn
	StateOff
)

// Switcher defines the interface for a toggleable switch
type Switcher interface {
	SwitchOn(gpio.CodeTransmitter) error
	SwitchOff(gpio.CodeTransmitter) error
	ToggleState(gpio.CodeTransmitter) error
}

// Outlet type definition
type Outlet struct {
	Identifier  string `yaml:"identifier" json:"identifier"`
	PulseLength int    `yaml:"pulse_length" json:"pulse_length"`
	Protocol    int    `yaml:"protocol" json:"protocol"`
	CodeOn      uint64 `yaml:"code_on" json:"code_on"`
	CodeOff     uint64 `yaml:"code_off" json:"code_off"`
	State       int    `yaml:"state" json:"state"`
}

// NewOutlet creates a new outlet
func NewOutlet(identifier string, gpioPin int, pulseLength int, protocol int, codeOn uint64, codeOff uint64) *Outlet {
	return &Outlet{
		Identifier:  identifier,
		PulseLength: pulseLength,
		Protocol:    protocol,
		CodeOn:      codeOn,
		CodeOff:     codeOff,
		State:       StateUnknown,
	}
}

func (o *Outlet) ToggleState(t gpio.CodeTransmitter) error {
	switch o.State {
	case StateOn:
		return o.SwitchOff(t)
	default:
		return o.SwitchOn(t)
	}
}

func (o *Outlet) SwitchOn(t gpio.CodeTransmitter) error {
	if err := o.sendCode(t, o.CodeOn); err != nil {
		return err
	}

	o.State = StateOn

	return nil
}

func (o *Outlet) SwitchOff(t gpio.CodeTransmitter) error {
	if err := o.sendCode(t, o.CodeOff); err != nil {
		return err
	}

	o.State = StateOff

	return nil
}

func (o *Outlet) sendCode(t gpio.CodeTransmitter, code uint64) error {
	return t.Transmit(code, o.Protocol, o.PulseLength)
}

// String returns the string representation of an Outlet
func (o *Outlet) String() string {
	return fmt.Sprintf("Outlet{Identifier: \"%s\", PulseLength: %d, Protocol: %d, CodeOn: %d, CodeOff: %d, State: %d}",
		o.Identifier, o.PulseLength, o.Protocol, o.CodeOn, o.CodeOff, o.State)
}

// UnmarshalYAML sets defaults on the raw Outlet before unmarshalling
func (o *Outlet) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawOutlet Outlet

	raw := rawOutlet{
		PulseLength: gpio.DefaultPulseLength,
		Protocol:    gpio.DefaultProtocol,
		State:       StateUnknown,
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*o = Outlet(raw)

	return nil
}
