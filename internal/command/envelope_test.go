package command

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name          string
		envelope      Envelope
		expected      Command
		expectedError error
	}{
		{"status command", Envelope{Type: StatusType}, &StatusCommand{}, nil},
		{"outlet command", Envelope{Type: OutletType, Data: rawMessage(`{"outletID":"foo","action":"toggle"}`)}, &OutletCommand{OutletID: "foo", Action: "toggle"}, nil},
		{"group command", Envelope{Type: GroupType, Data: rawMessage(`{"groupID":"foo","action":"on"}`)}, &GroupCommand{GroupID: "foo", Action: "on"}, nil},
		{"interval command", Envelope{Type: IntervalType, Data: rawMessage(`{"outletID":"foo","action":"create"}`)}, &IntervalCommand{OutletID: "foo", Action: "create"}, nil},
		// Error cases.
		{"unknown command", Envelope{Type: Type("unknown")}, nil, errors.New(`unknown command type "unknown"`)},
		{"invalid data", Envelope{Type: OutletType, Data: rawMessage(`{`)}, nil, errors.New(`failed to unmarshal "outlet" command data: unexpected end of JSON input`)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd, err := Unpack(test.envelope)
			if test.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, cmd)
			}
		})
	}
}

func rawMessage(s string) *json.RawMessage {
	raw := json.RawMessage(s)
	return &raw
}
