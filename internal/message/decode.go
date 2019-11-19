package message

import (
	"encoding/json"
)

// Decode decodes the contents of a message envelope into the correct type
func Decode(envelope Envelope) (interface{}, error) {
	switch envelope.Type {
	case OutletActionType:
		return decode(envelope.Data, &OutletAction{})
	case GroupActionType:
		return decode(envelope.Data, &GroupAction{})
	case IntervalActionType:
		return decode(envelope.Data, &IntervalAction{})
	}

	return &Unknown{}, nil
}

func decode(data *json.RawMessage, msg interface{}) (interface{}, error) {
	return msg, json.Unmarshal(*data, msg)
}
