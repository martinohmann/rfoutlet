package message

import (
	"encoding/json"

	"github.com/martinohmann/rfoutlet/internal/schedule"
)

const (
	outletActionType   = "outlet"
	groupActionType    = "group"
	intervalActionType = "interval"
)

// Dispatcher defines the interface for a message dispatcher
type Dispatcher interface {
	Dispatch(Envelope) error
}

// Unknown defines an unknown message
type Unknown struct{}

// OutletAction defines an outlet action message
type OutletAction struct {
	ID     string
	Action string
}

// GroupAction defines a group action message
type GroupAction struct {
	ID     string
	Action string
}

// IntervalAction defines an interval message
type IntervalAction struct {
	ID       string
	Interval schedule.Interval
	Action   string
}

// Envelope defines a message envelope which hold the message type and the raw
// json data of the message that gets unmarshalled into the correct type by
// Decode.
type Envelope struct {
	Type string
	Data *json.RawMessage
}
