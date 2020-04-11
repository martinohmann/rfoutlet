package config

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/imdario/mergo"
	"github.com/martinohmann/rfoutlet/internal/outlet"
)

const (
	// DefaultListenAddress defines the default address to listen on.
	DefaultListenAddress = ":3333"

	// DefaultTransmitPin defines the default gpio pin for transmitting rf codes.
	DefaultTransmitPin uint = 17

	// DefaultReceivePin defines the default gpio pin for receiving rf codes.
	DefaultReceivePin uint = 27

	// DefaultProtocol defines the default rf protocol.
	DefaultProtocol int = 1

	// DefaultPulseLength defines the default pulse length.
	DefaultPulseLength uint = 189
)

var DefaultConfig = Config{
	ListenAddress: DefaultListenAddress,
	ReceivePin:    DefaultReceivePin,
	TransmitPin:   DefaultTransmitPin,
}

type Config struct {
	ListenAddress string              `json:"listenAddress"`
	StateFile     string              `json:"stateFile"`
	ReceivePin    uint                `json:"receivePin"`
	TransmitPin   uint                `json:"transmitPin"`
	OutletGroups  []OutletGroupConfig `json:"outletGroups"`
}

type OutletGroupConfig struct {
	ID          string         `json:"id"`
	DisplayName string         `json:"displayName"`
	Outlets     []OutletConfig `outlets`
}

type OutletConfig struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	CodeOn      uint64 `json:"codeOn"`
	CodeOff     uint64 `json:"codeOff"`
	Protocol    int    `json:"protocol"`
	PulseLength uint   `json:"pulseLength"`
}

// BuildOutletGroups builds outlet groups from c.
func (c Config) BuildOutletGroups() []*outlet.Group {
	groups := make([]*outlet.Group, len(c.OutletGroups))

	for i, gc := range c.OutletGroups {
		outlets := make([]*outlet.Outlet, len(gc.Outlets))

		for j, oc := range gc.Outlets {
			outlets[j] = &outlet.Outlet{
				ID:          oc.ID,
				DisplayName: oc.DisplayName,
				CodeOn:      oc.CodeOn,
				CodeOff:     oc.CodeOff,
				Protocol:    oc.Protocol,
				PulseLength: oc.PulseLength,
			}
		}

		groups[i] = &outlet.Group{
			ID:          gc.ID,
			DisplayName: gc.DisplayName,
			Outlets:     outlets,
		}
	}

	return groups
}

func LoadWithDefaults(file string) (*Config, error) {
	config, err := Load(file)
	if err != nil {
		return nil, err
	}

	err = mergo.Merge(config, DefaultConfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Load loads the config from a file
func Load(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return LoadWithReader(f)
}

// LoadWithReader loads the config using reader
func LoadWithReader(r io.Reader) (*Config, error) {
	c, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = yaml.Unmarshal(c, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
