package config

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	yaml "gopkg.in/yaml.v2"
)

// DefaultListenAddress defines the default address to listen on
const DefaultListenAddress = ":3333"

// Config type definition
type Config struct {
	ListenAddress string             `yaml:"listen_address"`
	GpioPin       uint               `yaml:"gpio_pin"`
	StateFile     string             `yaml:"state_file"`
	GroupOrder    []string           `yaml:"group_order"`
	Groups        map[string]*Group  `yaml:"groups"`
	Outlets       map[string]*Outlet `yaml:"outlets"`
}

// UnmarshalYAML sets defaults on the raw Config before unmarshalling
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawConfig Config

	raw := rawConfig{
		ListenAddress: DefaultListenAddress,
		GpioPin:       gpio.DefaultTransmitPin,
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*c = Config(raw)

	return nil
}

// Group config type definition
type Group struct {
	Name    string   `yaml:"name"`
	Outlets []string `yaml:"outlets"`
}

// Outlet config type definition
type Outlet struct {
	Name        string `yaml:"name"`
	CodeOn      uint64 `yaml:"code_on"`
	CodeOff     uint64 `yaml:"code_off"`
	Protocol    int    `yaml:"protocol"`
	PulseLength uint   `yaml:"pulse_length"`
}

// UnmarshalYAML sets defaults on the raw Outlet before unmarshalling
func (o *Outlet) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawOutlet Outlet

	raw := rawOutlet{
		PulseLength: gpio.DefaultPulseLength,
		Protocol:    gpio.DefaultProtocol,
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	*o = Outlet(raw)

	return nil
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

	err = yaml.UnmarshalStrict(c, config)

	return config, err
}
