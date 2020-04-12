package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// Type is the type of a command.
type Type string

// Command types.
const (
	GroupType    Type = "group"
	IntervalType Type = "interval"
	OutletType   Type = "outlet"
	StatusType   Type = "status"
)

// ClientAwareCommand is aware of the websocket client.
type ClientAwareCommand interface {
	command.Command

	// SetClient sets the websocket client on the command. The client can be
	// used to send messages back to the client that issued the command.
	SetClient(client *Client)
}

// StatusCommand is sent by a connected client to retrieve the list of current
// outlet groups. This usually happens when the client first connects.
type StatusCommand struct {
	client *Client
}

func (c StatusCommand) Execute(context command.Context) (bool, error) {
	groups := context.GetGroups()

	msg, err := json.Marshal(groups)
	if err != nil {
		return false, err
	}

	c.client.Send(msg)

	return false, nil
}

func (c *StatusCommand) SetClient(client *Client) {
	c.client = client
}

// OutletCommand...
type OutletCommand struct {
	OutletID string `json:"id"`
	Action   string `json:"action"`
}

func (c OutletCommand) Execute(context command.Context) (bool, error) {
	outlet, ok := context.GetOutlet(c.OutletID)
	if !ok {
		return false, fmt.Errorf("outlet %q does not exist", c.OutletID)
	}

	if outlet.Schedule.Enabled() {
		return false, nil
	}

	err := context.Switch(outlet, getTargetState(outlet, c.Action))
	if err != nil {
		return false, err
	}

	return true, nil
}

// GroupCommand...
type GroupCommand struct {
	GroupID string `json:"id"`
	Action  string `json:"action"`
}

func (c GroupCommand) Execute(context command.Context) (bool, error) {
	group, ok := context.GetGroup(c.GroupID)
	if !ok {
		return false, fmt.Errorf("outlet group %q does not exist", c.GroupID)
	}

	var modified bool

	for _, outlet := range group.Outlets {
		if outlet.Schedule.Enabled() {
			continue
		}

		err := context.Switch(outlet, getTargetState(outlet, c.Action))
		if err != nil {
			return modified, err
		}

		modified = true
	}

	return modified, nil
}

// IntervalCommand...
type IntervalCommand struct {
	OutletID string            `json:"id"`
	Action   string            `json:"action"`
	Interval schedule.Interval `json:"interval"`
}

func (c IntervalCommand) Execute(context command.Context) (bool, error) {
	outlet, ok := context.GetOutlet(c.OutletID)
	if !ok {
		return false, fmt.Errorf("outlet %q does not exist", c.OutletID)
	}

	err := c.handle(outlet)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c IntervalCommand) handle(outlet *outlet.Outlet) error {
	switch c.Action {
	case "create":
		return outlet.Schedule.AddInterval(c.Interval)
	case "update":
		return outlet.Schedule.UpdateInterval(c.Interval)
	case "delete":
		return outlet.Schedule.DeleteInterval(c.Interval)
	default:
		return fmt.Errorf("invalid interval action %q", c.Action)
	}
}

func getTargetState(o *outlet.Outlet, action string) outlet.State {
	switch action {
	default:
		return outlet.StateOn
	case "off":
		return outlet.StateOff
	case "toggle":
		if o.GetState() == outlet.StateOn {
			return outlet.StateOff
		}

		return outlet.StateOn
	}
}

// Envelope defines a command envelope which hold the command type and the raw
// json data of the command that gets unmarshalled into the correct type by
// decodeCommand.
type Envelope struct {
	Type Type
	Data *json.RawMessage
}

// decodeCommand decodes the contents of a command envelope into the correct
// type.
func decodeCommand(envelope Envelope) (command.Command, error) {
	switch envelope.Type {
	case OutletType:
		return decode(envelope.Data, &OutletCommand{})
	case GroupType:
		return decode(envelope.Data, &GroupCommand{})
	case IntervalType:
		return decode(envelope.Data, &IntervalCommand{})
	case StatusType:
		return &StatusCommand{}, nil
	default:
		return nil, fmt.Errorf("unknown message type %q", envelope.Type)
	}
}

func decode(data *json.RawMessage, cmd command.Command) (command.Command, error) {
	return cmd, json.Unmarshal(*data, cmd)
}
