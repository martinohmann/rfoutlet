package command

import (
	"encoding/json"
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/outlet"
)

// Context is passed to every command.
type Context struct {
	// Registry contains all known outlets and outlet groups.
	Registry *outlet.Registry
	// Switcher can switch an outlet on or off.
	Switcher outlet.Switcher
}

type Command interface {
	// Execute executes the command. The passed in context contains the outlet
	// registry and switcher to be able to manipulate the state of outlets.
	// Must return true if a potential state change should be broadcasted to
	// all connected clients.
	Execute(context Context) (broadcast bool, err error)
}

// Sender can send a message.
type Sender interface {
	// Send sends a message.
	Send(msg []byte)
}

// SenderAwareCommand is aware the the command sender and can send messages
// back.
type SenderAwareCommand interface {
	Command

	// SetSender sets the sender on the command. The sender can be used to send
	// messages back.
	SetSender(sender Sender)
}

// Envelope defines a command envelope which hold the command type and the raw
// json data of the command that gets unmarshalled into the correct type by
// Decode.
type Envelope struct {
	Type Type
	Data *json.RawMessage
}

// Decode decodes the contents of a command envelope into the correct type.
func Decode(envelope Envelope) (Command, error) {
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

func decode(data *json.RawMessage, cmd Command) (Command, error) {
	return cmd, json.Unmarshal(*data, cmd)
}
