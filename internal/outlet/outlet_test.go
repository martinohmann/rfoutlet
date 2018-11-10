package outlet_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

var transmitter, _ = gpio.NewNullTransmitter()

func TestOutletSwitchOn(t *testing.T) {
	o := &outlet.Outlet{CodeOn: 1, CodeOff: 2, State: outlet.StateUnknown, Protocol: 1}

	err := o.SwitchOn(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOn, o.State)
}

func TestOutletSwitchOff(t *testing.T) {
	o := &outlet.Outlet{CodeOn: 1, CodeOff: 2, State: outlet.StateUnknown, Protocol: 1}

	err := o.SwitchOff(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOff, o.State)
}

func TestOutletToggleState(t *testing.T) {
	o := &outlet.Outlet{CodeOn: 1, CodeOff: 2, State: outlet.StateUnknown, Protocol: 1}

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

func TestUnmarshalError(t *testing.T) {
	o := &outlet.Outlet{}

	err := yaml.Unmarshal([]byte("[]"), o)

	assert.NotNil(t, err)
}
