package outlet_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	c, err := outlet.ReadConfig("testdata/valid-config.yml")

	if assert.NoError(t, err) && assert.Len(t, c.OutletGroups, 2) {
		assert.Equal(t, c.OutletGroups[0].Identifier, "Living Room")
	}
}

func TestReadMissingConfig(t *testing.T) {
	_, err := outlet.ReadConfig("testdata/nonexistent-config.yml")

	assert.EqualError(t, err, "open testdata/nonexistent-config.yml: no such file or directory")
}

func TestReadInvalidConfig(t *testing.T) {
	_, err := outlet.ReadConfig("testdata/invalid-config.yml")

	assert.EqualError(t, err, "yaml: unmarshal errors:\n  line 2: cannot unmarshal !!str `foo` into outlet.Config")
}
