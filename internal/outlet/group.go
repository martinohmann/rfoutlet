package outlet

import "fmt"

// OutletGroup type definition
type OutletGroup struct {
	Identifier string    `yaml:"identifier" json:"identifier"`
	Outlets    []*Outlet `yaml:"outlets" json:"outlets"`
}

// NewOutletGroup creates a new outlet group
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
