package message

import (
	"encoding/json"
)

// Decode decodes the contents of a message envelope into the correct type
func Decode(env Envelope) (interface{}, error) {
	switch env.Type {
	case outletActionType:
		var msg OutletAction
		err := json.Unmarshal(*env.Data, &msg)

		return msg, err
	case groupActionType:
		var msg GroupAction
		err := json.Unmarshal(*env.Data, &msg)

		return msg, err
	case intervalActionType:
		var msg IntervalAction
		err := json.Unmarshal(*env.Data, &msg)

		return msg, err
	}

	return Unknown{}, nil
}
