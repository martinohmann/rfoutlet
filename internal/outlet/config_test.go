package outlet_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
)

func TestOutletGroup(t *testing.T) {
	og := &outlet.OutletGroup{}

	c := &outlet.Config{
		OutletGroups: []*outlet.OutletGroup{og},
	}

	res, err := c.OutletGroup(0)

	assert.Nil(t, err)
	assert.Equal(t, og, res)

	res, err = c.OutletGroup(1)

	assert.Nil(t, res)
	assert.EqualError(t, err, "invalid offset 1")
}

func TestReadConfig(t *testing.T) {
	c, err := outlet.ReadConfig("../../dist/config.yml")

	assert.Nil(t, err)
	if assert.Len(t, c.OutletGroups, 2) {
		assert.Equal(t, c.OutletGroups[0].Identifier, "Living Room")
	}
}

func TestReadMissingConfig(t *testing.T) {
	_, err := outlet.ReadConfig("../../dist/nonexistent.yml")

	assert.EqualError(t, err, "open ../../dist/nonexistent.yml: no such file or directory")
}
