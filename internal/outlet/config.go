package outlet

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config type definition
type Config struct {
	OutletGroups []*OutletGroup `yaml:"outlet_groups" json:"outlet_groups"`
}

// OutletGroup returns the outlet group at given offset in the config
func (c *Config) OutletGroup(offset int) (*OutletGroup, error) {
	if offset >= 0 && len(c.OutletGroups) > offset {
		return c.OutletGroups[offset], nil
	}

	return nil, fmt.Errorf("invalid offset %d", offset)
}

// ReadConfig reads the outlet config from a yaml file
func ReadConfig(filename string) (*Config, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = yaml.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
