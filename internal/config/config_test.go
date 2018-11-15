package config_test

import (
	"testing"

	"github.com/Flaque/filet"
	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadWithReader(t *testing.T) {
	defer filet.CleanUp(t)

	tests := []struct {
		content string
		wantErr bool
		errMsg  string
		assert  func(t *testing.T, c *config.Config)
	}{
		{
			content: "",
			wantErr: false,
			assert: func(t *testing.T, c *config.Config) {
				assert.Len(t, c.Groups, 0)
			},
		},
		{
			content: "groups:\n  foo:\n    name: bar",
			wantErr: false,
			assert: func(t *testing.T, c *config.Config) {
				assert.Equal(t, "bar", c.Groups["foo"].Name)
			},
		},
		{
			content: "a\n- a\n[",
			wantErr: true,
			errMsg:  "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `a - a [` into config.Config",
		},
	}

	for _, tt := range tests {
		f := filet.TmpFile(t, "/tmp", tt.content)
		f.Seek(0, 0)

		c, err := config.LoadWithReader(f)

		if tt.wantErr {
			if assert.Error(t, err) {
				assert.Equal(t, tt.errMsg, err.Error())
			}
		} else if assert.NoError(t, err) {
			tt.assert(t, c)
		}
	}
}
