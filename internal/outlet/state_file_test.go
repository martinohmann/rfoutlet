package outlet

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateFile_ReadBack(t *testing.T) {
	sf := NewStateFile("testdata/valid_state.json")

	outlets := []*Outlet{{ID: "foo"}, {ID: "bar"}, {ID: "baz"}}

	expected := []*Outlet{
		{ID: "foo", State: StateOn, Schedule: schedule.NewWithIntervals([]schedule.Interval{
			{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Monday},
				From:     schedule.NewDayTime(0, 59),
				To:       schedule.NewDayTime(2, 1),
			},
		})},
		{ID: "bar", State: StateOn, Schedule: schedule.New()},
		{ID: "baz"},
	}

	require.NoError(t, sf.ReadBack(outlets))
	assert.Equal(t, expected, outlets)
}

func TestStateFile_ReadBack_Invalid(t *testing.T) {
	sf := NewStateFile("testdata/invalid_state.json")
	outlets := []*Outlet{}
	require.Error(t, sf.ReadBack(outlets))
}

func TestStateFile_WriteOut(t *testing.T) {
	f, err := ioutil.TempFile("", "rfoutlet-state-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	sf := NewStateFile(f.Name())

	outlets := []*Outlet{
		{ID: "foo", State: StateOn, Schedule: schedule.NewWithIntervals([]schedule.Interval{
			{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Monday},
				From:     schedule.NewDayTime(0, 59),
				To:       schedule.NewDayTime(2, 1),
			},
		})},
		{ID: "bar", State: StateOn, Schedule: schedule.New()},
		{ID: "baz"},
	}

	require.NoError(t, sf.WriteOut(outlets))

	buf, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"bar":{"state":1,"schedule":[]},"baz":{},"foo":{"state":1,"schedule":[{"id":"","enabled":true,"weekdays":[1],"from":{"hour":0,"minute":59},"to":{"hour":2,"minute":1}}]}}`

	assert.Equal(t, expected, string(buf))
}
