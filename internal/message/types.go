package message

import (
	"encoding/json"

	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// Type is the type of a message.
type Type string

const (
	GroupType    Type = "group"
	IntervalType Type = "interval"
	OutletType   Type = "outlet"
	StatusType   Type = "status"
)

// Dispatcher defines the interface for a message dispatcher
type Dispatcher interface {
	Dispatch(Envelope) error
}

// Message is the interface for a message.
type Message interface{}

// StatusMessage...
type StatusMessage struct{}

// OutletMessage...
type OutletMessage struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

// GroupMessage...
type GroupMessage struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

// IntervalMessage...
type IntervalMessage struct {
	ID       string            `json:"id"`
	Action   string            `json:"action"`
	Interval schedule.Interval `json:"interval"`
}

// Envelope defines a message envelope which hold the message type and the raw
// json data of the message that gets unmarshalled into the correct type by
// Decode.
type Envelope struct {
	Type Type
	Data *json.RawMessage
}
