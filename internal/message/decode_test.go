package message

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeMessage(t *testing.T) {
	tests := []struct {
		t             Type
		data          string
		expected      interface{}
		expectedError error
	}{
		{t: OutletType, data: `{"id": "foo","action":"on"}`, expected: &OutletMessage{}},
		{t: GroupType, data: `{"id": "bar","action":"toggle"}`, expected: &GroupMessage{}},
		{t: IntervalType, data: `{"id":"baz","action":"create","interval":{"ID":"foo"}}`, expected: &IntervalMessage{}},
		{t: StatusType, data: `{}`, expected: &StatusMessage{}},
		{t: "foo", expectedError: errors.New(`unknown message type "foo"`)},
	}

	for _, tt := range tests {
		data := json.RawMessage([]byte(tt.data))
		env := Envelope{
			Type: tt.t,
			Data: &data,
		}

		msg, err := Decode(env)

		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.IsType(t, tt.expected, msg)
		}
	}
}

func TestDecodeInvalidMessage(t *testing.T) {
	data := json.RawMessage([]byte(`{`))
	env := Envelope{
		Type: OutletType,
		Data: &data,
	}

	_, err := Decode(env)

	assert.Error(t, err)
}
