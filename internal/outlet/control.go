package outlet

import (
	"fmt"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

// Control type definition
type Control struct {
	config       *Config
	stateManager StateManager
	transmitter  gpio.CodeTransmitter
}

// NewControl create a new outlet control
func NewControl(config *Config, stateManager StateManager, transmitter gpio.CodeTransmitter) *Control {
	return &Control{config: config, stateManager: stateManager, transmitter: transmitter}
}

// OutletGroups returns all known outlet groups
func (c *Control) OutletGroups() []*OutletGroup {
	return c.config.OutletGroups
}

// OutletGroup returns the outlet group at given offset in the config
func (c *Control) OutletGroup(offset int) (*OutletGroup, error) {
	if offset >= 0 && len(c.config.OutletGroups) > offset {
		return c.config.OutletGroups[offset], nil
	}

	return nil, fmt.Errorf("invalid outlet group offset %d", offset)
}

// Outlet returns the outlet at given offset in given group in the config
func (c *Control) Outlet(groupId int, offset int) (*Outlet, error) {
	og, err := c.OutletGroup(groupId)
	if err != nil {
		return nil, err
	}

	if offset >= 0 && len(og.Outlets) > offset {
		return og.Outlets[offset], nil
	}

	return nil, fmt.Errorf("invalid outlet offset %d in group %d", offset, groupId)
}

// RestoreState restores the outlet state
func (c *Control) RestoreState() error {
	return c.stateManager.RestoreState(c)
}

// SaveState saves the outlet state
func (c *Control) SaveState() error {
	return c.stateManager.SaveState(c)
}

// SwitchOn switches switch on
func (c *Control) SwitchOn(s Switcher) error {
	return c.doSwitch(s.SwitchOn)
}

// SwitchOff switches switch off
func (c *Control) SwitchOff(s Switcher) error {
	return c.doSwitch(s.SwitchOff)
}

// ToggleState toggles switch state
func (c *Control) ToggleState(s Switcher) error {
	return c.doSwitch(s.ToggleState)
}

func (c *Control) doSwitch(f func(gpio.CodeTransmitter) error) error {
	if err := f(c.transmitter); err != nil {
		return err
	}

	return c.SaveState()
}
