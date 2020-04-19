package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/controller/commands"
)

// Type is the type of a command.
type Type string

// Command that are expected from websocket clients.
const (
	GroupType    Type = "group"
	IntervalType Type = "interval"
	OutletType   Type = "outlet"
	StatusType   Type = "status"
)

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
		return decode(envelope.Data, &commands.OutletCommand{})
	case GroupType:
		return decode(envelope.Data, &commands.GroupCommand{})
	case IntervalType:
		return decode(envelope.Data, &commands.IntervalCommand{})
	case StatusType:
		return &commands.StatusCommand{}, nil
	default:
		return nil, fmt.Errorf("unknown command type %q", envelope.Type)
	}
}

func decode(data *json.RawMessage, cmd command.Command) (command.Command, error) {
	return cmd, json.Unmarshal(*data, cmd)
}
