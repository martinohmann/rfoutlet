package outlet

import (
	"fmt"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

// OutletGroup type definition
type OutletGroup struct {
	Identifier string    `yaml:"identifier" json:"identifier"`
	Outlets    []*Outlet `yaml:"outlets" json:"outlets"`
}

// NewOutletGroup creates a new outlet group
func NewOutletGroup(identifier string) *OutletGroup {
	return &OutletGroup{Identifier: identifier}
}

// AddOutlet adds an outlet to the group
func (og *OutletGroup) AddOutlet(outlet *Outlet) {
	og.Outlets = append(og.Outlets, outlet)
}

// Outlet returns the outlet with the given offset in the group
func (og *OutletGroup) Outlet(offset int) (*Outlet, error) {
	if offset >= 0 && len(og.Outlets) > offset {
		return og.Outlets[offset], nil
	}

	return nil, fmt.Errorf("invalid offset %d", offset)
}

// ToggleState toggles the state of all outlets of the group
func (og *OutletGroup) ToggleState(t gpio.CodeTransmitter) error {
	for _, o := range og.Outlets {
		if err := o.ToggleState(t); err != nil {
			return err
		}
	}

	return nil
}

// SwitchOn switches all outlets of the group on
func (og *OutletGroup) SwitchOn(t gpio.CodeTransmitter) error {
	for _, o := range og.Outlets {
		if err := o.SwitchOn(t); err != nil {
			return err
		}
	}

	return nil
}

// SwitchOff switches all outlets of the group off
func (og *OutletGroup) SwitchOff(t gpio.CodeTransmitter) error {
	for _, o := range og.Outlets {
		if err := o.SwitchOff(t); err != nil {
			return err
		}
	}

	return nil
}

// String returns the string representation of an OutletGroup
func (og *OutletGroup) String() string {
	return fmt.Sprintf("OutletGroup{Identifier: \"%s\"}", og.Identifier)
}
