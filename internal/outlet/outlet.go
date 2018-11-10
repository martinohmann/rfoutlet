package outlet

import (
	"log"
	"os"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "outlet: ", log.LstdFlags|log.Lshortfile)
}

// Switcher defines the interface for a toggleable switch
type Switcher interface {
	SwitchOn(gpio.CodeTransmitter) error
	SwitchOff(gpio.CodeTransmitter) error
	ToggleState(gpio.CodeTransmitter) error
}

// Outlet type definition
type Outlet struct {
	Identifier  string `yaml:"identifier" json:"identifier"`
	PulseLength uint   `yaml:"pulse_length" json:"-"`
	Protocol    int    `yaml:"protocol" json:"-"`
	CodeOn      uint64 `yaml:"code_on" json:"-"`
	CodeOff     uint64 `yaml:"code_off" json:"-"`
	State       State  `yaml:"state" json:"state"`
}

// ToggleState toggles the state of the outlet
func (o *Outlet) ToggleState(t gpio.CodeTransmitter) error {
	switch o.State {
	case StateOn:
		return o.SwitchOff(t)
	default:
		return o.SwitchOn(t)
	}
}

// SwitchOn switches the outlet on
func (o *Outlet) SwitchOn(t gpio.CodeTransmitter) error {
	if err := o.sendCode(t, o.CodeOn); err != nil {
		return err
	}

	o.State = StateOn

	return nil
}

// SwitchOff switches the outlet off
func (o *Outlet) SwitchOff(t gpio.CodeTransmitter) error {
	if err := o.sendCode(t, o.CodeOff); err != nil {
		return err
	}

	o.State = StateOff

	return nil
}

func (o *Outlet) sendCode(t gpio.CodeTransmitter, code uint64) error {
	logger.Printf("transmitting code=%d pulseLength=%d protocol=%d\n", code, o.PulseLength, o.Protocol)

	return t.Transmit(code, o.Protocol, o.PulseLength)
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
