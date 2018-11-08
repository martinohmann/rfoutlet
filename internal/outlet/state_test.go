package outlet_test

import (
	"io/ioutil"
	"testing"

	"github.com/Flaque/filet"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var stateTestConfig = &outlet.Config{
	OutletGroups: []*outlet.OutletGroup{
		&outlet.OutletGroup{
			Outlets: []*outlet.Outlet{
				&outlet.Outlet{State: 2},
			},
		},
		&outlet.OutletGroup{
			Outlets: []*outlet.Outlet{
				&outlet.Outlet{State: 1},
				&outlet.Outlet{State: 0},
			},
		},
	},
}

func TestSaveStateEmptyFile(t *testing.T) {
	defer filet.CleanUp(t)

	f := filet.TmpFile(t, "/tmp", "")

	testSaveState(t, f, stateTestConfig)
}

func TestSaveStateOverwriteFile(t *testing.T) {
	defer filet.CleanUp(t)

	f := filet.TmpFile(t, "/tmp", "foo")

	testSaveState(t, f, stateTestConfig)
}

func testSaveState(t *testing.T, f afero.File, config *outlet.Config) {
	sm := outlet.NewStateManager(f)
	c := outlet.NewControl(config, sm, transmitter)

	err := c.SaveState()
	assert.Nil(t, err)

	f.Seek(0, 0)

	b, err := ioutil.ReadAll(f)
	assert.Nil(t, err)

	assert.Equal(t, "[{\"outlet\":0,\"group\":0,\"state\":2},{\"outlet\":0,\"group\":1,\"state\":1},{\"outlet\":1,\"group\":1,\"state\":0}]\n", string(b))
}

func TestRestoreState(t *testing.T) {
	defer filet.CleanUp(t)

	tests := []struct {
		configProvider func() *outlet.Config
		fileContents   string
		wantErr        bool
		errMsg         string
		assertFunc     func(*testing.T, *outlet.Config)
	}{
		{
			configProvider: func() *outlet.Config {
				return &outlet.Config{
					OutletGroups: []*outlet.OutletGroup{
						&outlet.OutletGroup{
							Outlets: []*outlet.Outlet{
								&outlet.Outlet{},
							},
						},
						&outlet.OutletGroup{
							Outlets: []*outlet.Outlet{
								&outlet.Outlet{},
								&outlet.Outlet{},
							},
						},
					},
				}
			},
			fileContents: "[{\"outlet\":0,\"group\":0,\"state\":2},{\"outlet\":0,\"group\":1,\"state\":1},{\"outlet\":1,\"group\":1,\"state\":0}]\n",
			assertFunc: func(t *testing.T, config *outlet.Config) {
				assert.Equal(t, outlet.StateOff, config.OutletGroups[0].Outlets[0].State)
				assert.Equal(t, outlet.StateOn, config.OutletGroups[1].Outlets[0].State)
				assert.Equal(t, outlet.StateUnknown, config.OutletGroups[1].Outlets[1].State)
			},
		},
		{
			configProvider: func() *outlet.Config {
				return &outlet.Config{}
			},
			fileContents: "[{\"Outlet\":0,\"Group\":0,\"SwitchState\":2}]\n",
			wantErr:      true,
			errMsg:       "invalid outlet group offset 0",
		},
		{
			configProvider: func() *outlet.Config {
				return &outlet.Config{
					OutletGroups: []*outlet.OutletGroup{
						&outlet.OutletGroup{
							Outlets: []*outlet.Outlet{
								&outlet.Outlet{},
							},
						},
					},
				}
			},
			fileContents: "[{\"outlet\":0,\"group\":0,\"state\":2}, {\"outlet\":1,\"group\":0,\"state\":1}]\n",
			wantErr:      true,
			errMsg:       "invalid outlet offset 1 in group 0",
		},
		{
			configProvider: func() *outlet.Config {
				return &outlet.Config{}
			},
			fileContents: "{",
			wantErr:      true,
			errMsg:       "unexpected end of JSON input",
		},
		{
			configProvider: func() *outlet.Config {
				return &outlet.Config{}
			},
			fileContents: "",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		f := filet.TmpFile(t, "/tmp", tt.fileContents)
		config := tt.configProvider()
		sm := outlet.NewStateManager(f)
		c := outlet.NewControl(config, sm, transmitter)

		err := c.RestoreState()

		if tt.wantErr {
			assert.NotNil(t, err)
			assert.EqualError(t, err, tt.errMsg)
		} else {
			assert.Nil(t, err)

			if tt.assertFunc != nil {
				tt.assertFunc(t, config)
			}
		}
	}
}
