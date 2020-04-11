package control

import (
	"encoding/json"
	"testing"

	"github.com/martinohmann/rfoutlet/internal/message"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createControl(r *outlet.Registry) *Control {
	return New(r, outlet.NewSwitch(gpio.NewDiscardingTransmitter()), NewHub())
}

func TestSwitch(t *testing.T) {
	c := createControl(outlet.NewRegistry())
	o := &outlet.Outlet{ID: "foo", Protocol: 1}

	assert.NoError(t, c.Switch(o, outlet.StateOn))
}

func TestDispatch(t *testing.T) {
	r := outlet.NewRegistry()

	err := r.RegisterGroups(&outlet.Group{
		ID:          "group",
		DisplayName: "Group",
		Outlets: []*outlet.Outlet{
			{ID: "foo", DisplayName: "Foo", Protocol: 1, PulseLength: 1},
		},
	})
	require.NoError(t, err)

	c := createControl(r)

	data := json.RawMessage([]byte(`{"id":"foo","action":"on"}`))
	env := message.Envelope{
		Type: message.OutletType,
		Data: &data,
	}

	assert.NoError(t, c.Dispatch(env))
}

func TestDispatchError(t *testing.T) {
	c := createControl(outlet.NewRegistry())

	env := message.Envelope{
		Type: message.OutletType,
		Data: &json.RawMessage{},
	}

	assert.Error(t, c.Dispatch(env))
}
