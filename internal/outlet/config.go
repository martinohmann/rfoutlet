package outlet

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

const DefaultConfigFilename = "config.yml"

// Config type definition
type Config struct {
	OutletGroups []*OutletGroup `yaml:"outlet_groups" json:"outlet_groups"`
}

// Print prints the current config values
func (c *Config) Print() {
	for _, og := range c.OutletGroups {
		fmt.Printf("%s\n", og)

		for _, o := range og.Outlets {
			fmt.Printf("  %s\n", o)
		}
	}
}

func (c *Config) OutletGroup(offset int) (*OutletGroup, error) {
	if offset >= 0 && len(c.OutletGroups) > offset {
		return c.OutletGroups[offset], nil
	}
	return nil, fmt.Errorf("invalid offset %d", offset)
}

// ReadConfig reads the outlet config from a yaml file. Will panic if the file
// is not readable or if it contains invalid yaml
func ReadConfig(filename string) *Config {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	config := &Config{}

	err = yaml.Unmarshal(contents, config)
	if err != nil {
		panic(err)
	}

	return config
}
