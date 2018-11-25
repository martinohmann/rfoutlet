package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeMessage(t *testing.T) {
	tests := []struct {
		t        string
		data     string
		expected interface{}
	}{
		{t: outletActionType, data: `{"id": "foo","action":"on"}`, expected: OutletAction{}},
		{t: groupActionType, data: `{"id": "bar","action":"toggle"}`, expected: GroupAction{}},
		{t: intervalActionType, data: `{"id":"baz","action":"create","interval":{"ID":"foo"}}`, expected: IntervalAction{}},
		{t: "foo", expected: Unknown{}},
	}

	for _, tt := range tests {
		data := json.RawMessage([]byte(tt.data))
		env := Envelope{
			Type: tt.t,
			Data: &data,
		}

		msg, err := Decode(env)

		assert.NoError(t, err)
		assert.IsType(t, tt.expected, msg)
	}
}

func TestDecodeInvalidMessage(t *testing.T) {
	data := json.RawMessage([]byte(`{`))
	env := Envelope{
		Type: outletActionType,
		Data: &data,
	}

	_, err := Decode(env)

	assert.Error(t, err)
}
