package outlet_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
)

func TestOutlet(t *testing.T) {
	o := &outlet.Outlet{Protocol: 1}

	og := &outlet.OutletGroup{
		Outlets: []*outlet.Outlet{o},
	}

	res, err := og.Outlet(0)

	assert.Nil(t, err)
	assert.Equal(t, o, res)

	res, err = og.Outlet(1)

	assert.Nil(t, res)
	assert.EqualError(t, err, "invalid offset 1")
}

func TestOutputGroupSwitchOn(t *testing.T) {
	o := &outlet.Outlet{State: outlet.StateOff, Protocol: 1}

	og := &outlet.OutletGroup{
		Outlets: []*outlet.Outlet{o},
	}

	err := og.SwitchOn(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOn, o.State)
}

func TestOutputGroupSwitchOff(t *testing.T) {
	o := &outlet.Outlet{State: outlet.StateOn, Protocol: 1}

	og := &outlet.OutletGroup{
		Outlets: []*outlet.Outlet{o},
	}

	err := og.SwitchOff(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOff, o.State)
}

func TestOutputGroupToggleState(t *testing.T) {
	o := &outlet.Outlet{State: outlet.StateOff, Protocol: 1}

	og := &outlet.OutletGroup{
		Outlets: []*outlet.Outlet{o},
	}

	err := og.ToggleState(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOn, o.State)

	err = og.ToggleState(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOff, o.State)
}

func TestOutputGroupToggleStateInvalidOutletProtocol(t *testing.T) {
	o := &outlet.Outlet{State: outlet.StateOff, Protocol: 9999}

	og := &outlet.OutletGroup{
		Outlets: []*outlet.Outlet{o},
	}

	err := og.ToggleState(transmitter)

	assert.NotNil(t, err)
}
