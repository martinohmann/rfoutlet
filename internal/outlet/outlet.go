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

const DefaultPulseLength = 189

// Switcher defines the interface for a toggleable switch
type Switcher interface {
	SwitchOn() error
	SwitchOff() error
	ToggleState() error
}

// Outlet type definition
type Outlet struct {
	Identifier  string `yaml:"identifier" json:"identifier"`
	GpioPin     int    `yaml:"gpio_pin" json:"gpio_pin"`
	PulseLength int    `yaml:"pulse_length" json:"pulse_length"`
	CodeOn      int    `yaml:"code_on" json:"code_on"`
	CodeOff     int    `yaml:"code_off" json:"code_off"`
	State       int    `yaml:"state" json:"state"`
	transmit    gpio.TransmitFunc
}

// NewOutlet creates a new outlet
func NewOutlet(identifier string, gpioPin int, pulseLength int, codeOn int, codeOff int) *Outlet {
	return &Outlet{
		Identifier:  identifier,
		GpioPin:     gpioPin,
		PulseLength: pulseLength,
		CodeOn:      codeOn,
		CodeOff:     codeOff,
		State:       StateUnknown,
		transmit:    gpio.Transmit,
	}
}

func (o *Outlet) ToggleState() error {
	switch o.State {
	case StateOn:
		return o.SwitchOff()
	default:
		return o.SwitchOn()
	}
}

func (o *Outlet) SwitchOn() error {
	if err := o.transmit(o.CodeOn, o.GpioPin, o.PulseLength); err != nil {
		return err
	}

	o.State = StateOn

	return nil
}

func (o *Outlet) SwitchOff() error {
	if err := o.transmit(o.CodeOff, o.GpioPin, o.PulseLength); err != nil {
		return err
	}

	o.State = StateOff

	return nil
}

func (o *Outlet) SetTransmitFunc(f gpio.TransmitFunc) {
	o.transmit = f
}

// String returns the string representation of an Outlet
func (o *Outlet) String() string {
	return fmt.Sprintf("Outlet{Identifier: \"%s\", GpioPin: %d, PulseLength: %d, CodeOn: %d, CodeOff: %d, State: %d}",
		o.Identifier, o.GpioPin, o.PulseLength, o.CodeOn, o.CodeOff, o.State)
}

// UnmarshalYAML sets defaults on the raw Outlet before unmarshalling
func (o *Outlet) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawOutlet Outlet

	raw := rawOutlet{
		PulseLength: DefaultPulseLength,
		transmit:    gpio.Transmit,
		State:       StateUnknown,
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*o = Outlet(raw)

	return nil
}
