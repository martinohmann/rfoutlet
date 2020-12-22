package command

import (
	"encoding/json"
	"fmt"
)

// Envelope defines a command envelope which hold the command type and the raw
// json data of the command.
type Envelope struct {
	Type Type
	Data *json.RawMessage
}

// Unpack unpacks the contents of a command envelope into the correct
// type.
func Unpack(envelope Envelope) (cmd Command, err error) {
	switch envelope.Type {
	case OutletType:
		cmd = &OutletCommand{}
	case GroupType:
		cmd = &GroupCommand{}
	case IntervalType:
		cmd = &IntervalCommand{}
	case StatusType:
		cmd = &StatusCommand{}
	default:
		return nil, fmt.Errorf("unknown command type %q", envelope.Type)
	}

	if envelope.Data != nil {
		err = json.Unmarshal(*envelope.Data, cmd)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal %q command data: %w", envelope.Type, err)
		}
	}

	return cmd, nil
}
