package state

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadInvalid(t *testing.T) {
	_, err := Load("testdata/invalid.json")
	assert.Error(t, err)
}

func TestLoadNonexistent(t *testing.T) {
	_, err := Load("testdata/idonotexist.json")
	assert.Error(t, err)
}

func TestLoadValid(t *testing.T) {
	s, err := Load("testdata/valid.json")
	require.NoError(t, err)

	expected := State{
		"foo": OutletState{
			State: outlet.StateOn,
			Schedule: schedule.NewWithIntervals([]schedule.Interval{
				{
					Enabled: true,
					From: schedule.DayTime{
						Hour:   0,
						Minute: 59,
					},
					To: schedule.DayTime{
						Hour:   2,
						Minute: 1,
					},
					Weekdays: []time.Weekday{time.Monday},
				},
			}),
		},
	}

	assert.Equal(t, expected, s)
}

func TestLoadWithReader(t *testing.T) {
	r := strings.NewReader(`{"foo":{"state":1}}`)
	s, err := LoadWithReader(r)
	require.NoError(t, err)

	outletState, ok := s["foo"]
	require.True(t, ok)
	assert.Equal(t, outlet.StateOn, outletState.State)
}

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 1, fmt.Errorf("error")
}

func TestLoadWithBadReader(t *testing.T) {
	_, err := LoadWithReader(errorReader{})
	assert.Error(t, err)
}

func TestSaveInNonexistentDir(t *testing.T) {
	assert.Error(t, Save("testdata/thisdoesnotexist/json", State{}))
}

type errorWriter struct{}

func (errorWriter) Write(p []byte) (n int, err error) {
	return 1, fmt.Errorf("error")
}

func TestSaveWithWriter(t *testing.T) {
	var buf []byte

	w := bytes.NewBuffer(buf)

	assert.NoError(t, SaveWithWriter(w, State{}))
	assert.Equal(t, "{}\n", w.String())
}

func TestSaveWithBadWriter(t *testing.T) {
	assert.Error(t, SaveWithWriter(errorWriter{}, State{}))
}

func TestCollect(t *testing.T) {
	schedule := schedule.New()
	outlets := []*outlet.Outlet{{ID: "foo", State: outlet.StateOn, Schedule: schedule}}

	s := Collect(outlets)

	outletState, ok := s["foo"]
	require.True(t, ok)
	assert.Equal(t, schedule, outletState.Schedule)
	assert.Equal(t, outlet.StateOn, outletState.State)
}

func TestApply(t *testing.T) {
	outlets := []*outlet.Outlet{{ID: "foo", State: outlet.StateOn}}
	sch := schedule.New()

	s := State{
		"foo": OutletState{
			State:    outlet.StateOff,
			Schedule: sch,
		},
	}

	s.Apply(outlets)

	assert.Equal(t, outlet.StateOff, outlets[0].State)
	assert.Equal(t, sch, outlets[0].Schedule)
}
