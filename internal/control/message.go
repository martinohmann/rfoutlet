package control

import (
	"encoding/json"
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/schedule"
)

const (
	statusMessageType   = "status"
	outletMessageType   = "outlet"
	groupMessageType    = "group"
	intervalMessageType = "interval"
)

type outletMessage struct {
	ID     string
	Action string
}

type groupMessage struct {
	ID     string
	Action string
}

type intervalMessage struct {
	ID       string
	Interval schedule.Interval
	Action   string
}

type messageEnvelope struct {
	Type string
	Data *json.RawMessage
}

func decodeMessage(env messageEnvelope) (interface{}, error) {
	switch env.Type {
	case statusMessageType:
		return nil, nil
	case outletMessageType:
		var msg outletMessage
		err := json.Unmarshal(*env.Data, &msg)

		return msg, err
	case groupMessageType:
		var msg groupMessage
		err := json.Unmarshal(*env.Data, &msg)

		return msg, err
	case intervalMessageType:
		var msg intervalMessage
		err := json.Unmarshal(*env.Data, &msg)

		return msg, err
	}

	return nil, fmt.Errorf("unknown message type: %q", env.Type)
}
