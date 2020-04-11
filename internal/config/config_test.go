package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	c, err := Load("testdata/full.yaml")

	require.NoError(t, err)
	assert.Equal(t, "0.0.0.0:1234", c.ListenAddress)
	require.Len(t, c.OutletGroups, 2)
	assert.Len(t, c.OutletGroups[0].Outlets, 2)
	assert.Len(t, c.OutletGroups[1].Outlets, 1)
}

func TestLoadWithDefaults(t *testing.T) {
	c, err := LoadWithDefaults("testdata/partial.yaml")

	require.NoError(t, err)
	assert.Equal(t, DefaultConfig.ListenAddress, c.ListenAddress)
	assert.Empty(t, c.StateFile)
	assert.Equal(t, uint(42), c.ReceivePin)
	assert.Equal(t, DefaultConfig.TransmitPin, c.TransmitPin)
	require.Len(t, c.OutletGroups, 2)
}

func TestLoadInvalid(t *testing.T) {
	_, err := Load("testdata/invalid.yml")
	assert.Error(t, err)
}

func TestLoadNonexistent(t *testing.T) {
	_, err := Load("testdata/idonotexist.yml")
	assert.Error(t, err)
}

func TestLoadWithReader(t *testing.T) {
	cfg := `
outletGroups:
  - id: foo
    displayName: Foo`

	r := strings.NewReader(cfg)
	c, err := LoadWithReader(r)
	require.NoError(t, err)
	require.Len(t, c.OutletGroups, 1)
	assert.Equal(t, "foo", c.OutletGroups[0].ID)
	assert.Equal(t, "Foo", c.OutletGroups[0].DisplayName)
}

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 1, fmt.Errorf("error")
}

func TestLoadWithBadReader(t *testing.T) {
	_, err := LoadWithReader(errorReader{})
	assert.Error(t, err)
}

func TestConfig_BuildOutletGroups(t *testing.T) {
	config := Config{
		OutletGroups: []OutletGroupConfig{
			{
				ID:          "foo",
				DisplayName: "Foo",
				Outlets: []OutletConfig{
					{
						ID:      "bar",
						CodeOn:  1,
						CodeOff: 2,
					},
				},
			},
		},
	}

	expected := []*outlet.Group{
		{
			ID:          "foo",
			DisplayName: "Foo",
			Outlets: []*outlet.Outlet{
				{
					ID:      "bar",
					CodeOn:  1,
					CodeOff: 2,
				},
			},
		},
	}

	assert.Equal(t, expected, config.BuildOutletGroups())
}
