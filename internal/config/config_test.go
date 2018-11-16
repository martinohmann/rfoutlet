package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	c, err := config.Load("testdata/full.yml")

	assert.NoError(t, err)
	assert.Len(t, c.GroupOrder, 1)
	assert.Len(t, c.Groups, 1)
	assert.Len(t, c.Outlets, 2)
}

func TestLoadInvalid(t *testing.T) {
	_, err := config.Load("testdata/invalid.yml")
	assert.Error(t, err)
}

func TestLoadNonexistent(t *testing.T) {
	_, err := config.Load("testdata/idonotexist.yml")
	assert.Error(t, err)
}

func TestLoadWithReader(t *testing.T) {
	cfg := `
groups:
  foo:
    name: Foo`

	r := strings.NewReader(cfg)
	c, err := config.LoadWithReader(r)
	assert.NoError(t, err)
	assert.Equal(t, "Foo", c.Groups["foo"].Name)
}

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 1, fmt.Errorf("error")
}

func TestLoadWithBadReader(t *testing.T) {
	_, err := config.LoadWithReader(errorReader{})
	assert.Error(t, err)
}

var errorUnmarshal = func(interface{}) error {
	return fmt.Errorf("error")
}

func TestBadConfigUnmarshalYAML(t *testing.T) {
	c := &config.Config{}

	assert.Error(t, c.UnmarshalYAML(errorUnmarshal))
}

func TestBadOutletUnmarshalYAML(t *testing.T) {
	o := &config.Outlet{}

	assert.Error(t, o.UnmarshalYAML(errorUnmarshal))
}
