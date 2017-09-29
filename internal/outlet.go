package internal

import "fmt"

const (
	StateUnknown = iota
	StateOn
	StateOff
)

const DefaultPulseLength = 189

// Outlet type definition
type Outlet struct {
	Identifier  string `yaml:"identifier"`
	PulseLength int    `yaml:"pulse_length"`
	CodeOn      int    `yaml:"code_on"`
	CodeOff     int    `yaml:"code_off"`
	State       int    `yaml:"state"`
}

// OutletGroup type definition
type OutletGroup struct {
	Identifier string    `yaml:"identifier"`
	Outlets    []*Outlet `yaml:"outlets"`
}

// NewOutletGroup creates a new instance of the Outlet struct
func NewOutlet(identifier string, pulseLength int, codeOn int, codeOff int) *Outlet {
	return &Outlet{
		Identifier:  identifier,
		PulseLength: pulseLength,
		CodeOn:      codeOn,
		CodeOff:     codeOff,
		State:       StateUnknown,
	}
}

func (o *Outlet) ToggleState() error {
	switch o.State {
	case StateOn:
		return o.SwitchOff()
	case StateOff:
		return o.SwitchOn()
	default:
		return o.SwitchOn()
	}
}

func (o *Outlet) SwitchOn() error {
	if err := Transmit(o.CodeOn, o.PulseLength); err != nil {
		return err
	}
	o.State = StateOn
	return nil
}

func (o *Outlet) SwitchOff() error {
	if err := Transmit(o.CodeOff, o.PulseLength); err != nil {
		return err
	}
	o.State = StateOn
	return nil
}

// String returns the string representation of an Outlet
func (o *Outlet) String() string {
	return fmt.Sprintf("Outlet{Identifier: \"%s\", PulseLength: %d, CodeOn: %d, CodeOff: %d, State: %d}",
		o.Identifier, o.PulseLength, o.CodeOn, o.CodeOff, o.State)
}

// UnmarshalYAML sets defaults on the raw Outlet before unmarshalling
func (o *Outlet) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawOutlet Outlet

	raw := rawOutlet{
		PulseLength: DefaultPulseLength,
		State:       StateUnknown,
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*o = Outlet(raw)
	return nil
}

// NewOutletGroup creates a new instance of the OutletGroup struct
func NewOutletGroup(identifier string) *OutletGroup {
	return &OutletGroup{Identifier: identifier}
}

func (og *OutletGroup) AddOutlet(outlet *Outlet) {
	og.Outlets = append(og.Outlets, outlet)
}

func (og *OutletGroup) Outlet(offset int) (*Outlet, error) {
	if offset >= 0 && len(og.Outlets) > offset {
		return og.Outlets[offset], nil
	}
	return nil, fmt.Errorf("invalid offset %d", offset)
}

func (og *OutletGroup) ToggleState() error {
	for _, o := range og.Outlets {
		if err := o.ToggleState(); err != nil {
			return err
		}
	}
	return nil
}

func (og *OutletGroup) SwitchOn() error {
	for _, o := range og.Outlets {
		if err := o.SwitchOn(); err != nil {
			return err
		}
	}
	return nil
}

func (og *OutletGroup) SwitchOff() error {
	for _, o := range og.Outlets {
		if err := o.SwitchOff(); err != nil {
			return err
		}
	}
	return nil
}

// String returns the string representation of an OutletGroup
func (og *OutletGroup) String() string {
	return fmt.Sprintf("OutletGroup{Identifier: \"%s\"}", og.Identifier)
}
