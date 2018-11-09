package outlet_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	c, err := outlet.ReadConfig("../../example/config.yml")

	assert.Nil(t, err)
	if assert.Len(t, c.OutletGroups, 2) {
		assert.Equal(t, c.OutletGroups[0].Identifier, "Living Room")
	}
}

func TestReadMissingConfig(t *testing.T) {
	_, err := outlet.ReadConfig("../../example/nonexistent.yml")

	assert.EqualError(t, err, "open ../../example/nonexistent.yml: no such file or directory")
}
