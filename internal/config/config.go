package config

import (
	"io"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config type definition
type Config struct {
	ListenAddress string             `yaml:"listen_address,omitempty"`
	GpioPin       uint               `yaml:"gpio_pin,omitempty"`
	StateFile     string             `yaml:"state_file,omitempty"`
	GroupOrder    []string           `yaml:",omitempty"`
	Groups        map[string]*Group  `yaml:",omitempty"`
	Outlets       map[string]*Outlet `yaml:",omitempty"`
}

// Group config type definition
type Group struct {
	Name    string   `yaml:",omitempty"`
	Outlets []string `yaml:",omitempty"`
}

// Oultet config type definition
type Outlet struct {
	Name        string `yaml:",omitempty"`
	CodeOn      uint64 `yaml:"code_on,omitempty"`
	CodeOff     uint64 `yaml:"code_off,omitempty"`
	Protocol    int    `yaml:",omitempty"`
	PulseLength uint   `yaml:"pulse_length,omitempty"`
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
