package outlet_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
)

type mockTransmitter struct{}

func (t *mockTransmitter) Transmit(code uint64, protocol int, pulseLength int) error { return nil }
func (t *mockTransmitter) Close() error                                              { return nil }

var transmitter = &mockTransmitter{}

func TestNewOutlet(t *testing.T) {
	o := outlet.NewOutlet("foo", 1, 2, 3, 4)

	assert.Equal(t, "foo", o.Identifier)
	assert.Equal(t, 1, o.PulseLength)
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
