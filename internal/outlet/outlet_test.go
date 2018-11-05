package outlet_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

var transmitter, _ = gpio.NewNullTransmitter()

func TestNewOutlet(t *testing.T) {
	o := outlet.NewOutlet("foo", 1, 2, 3, 4)

	assert.Equal(t, "foo", o.Identifier)
	assert.Equal(t, uint(1), o.PulseLength)
	assert.Equal(t, 2, o.Protocol)
	assert.Equal(t, uint64(3), o.CodeOn)
	assert.Equal(t, uint64(4), o.CodeOff)
}

func TestOutletSwitchOn(t *testing.T) {
	o := &outlet.Outlet{CodeOn: 1, CodeOff: 2, State: outlet.StateUnknown}

	err := o.SwitchOn(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOn, o.State)
}

func TestOutletSwitchOff(t *testing.T) {
	o := &outlet.Outlet{CodeOn: 1, CodeOff: 2, State: outlet.StateUnknown}

	err := o.SwitchOff(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOff, o.State)
}

func TestOutletToggleState(t *testing.T) {
	o := &outlet.Outlet{CodeOn: 1, CodeOff: 2, State: outlet.StateUnknown}

	err := o.ToggleState(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOn, o.State)

	err = o.ToggleState(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOff, o.State)
}

func TestUnmarshalDefaults(t *testing.T) {
	o := &outlet.Outlet{}

	err := yaml.Unmarshal([]byte("{}"), o)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateUnknown, o.State)
	assert.Equal(t, gpio.DefaultPulseLength, o.PulseLength)
	assert.Equal(t, gpio.DefaultProtocol, o.Protocol)
}
