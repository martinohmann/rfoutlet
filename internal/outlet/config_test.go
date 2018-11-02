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
