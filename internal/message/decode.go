package message

import (
	"encoding/json"
	"fmt"
)

// Decode decodes the contents of a message envelope into the correct type
func Decode(envelope Envelope) (interface{}, error) {
	switch envelope.Type {
	case OutletType:
		return decode(envelope.Data, &OutletMessage{})
	case GroupType:
		return decode(envelope.Data, &GroupMessage{})
	case IntervalType:
		return decode(envelope.Data, &IntervalMessage{})
	case StatusType:
		return &StatusMessage{}, nil
	default:
		return nil, fmt.Errorf("unknown message type %q", envelope.Type)
	}
}

func decode(data *json.RawMessage, msg interface{}) (interface{}, error) {
	return msg, json.Unmarshal(*data, msg)
}
