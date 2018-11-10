package outlet

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config type definition
type Config struct {
	OutletGroups []*OutletGroup `yaml:"outlet_groups" json:"outlet_groups"`
}

// ReadConfig reads the outlet config from a yaml file
func ReadConfig(filename string) (*Config, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	if err = yaml.Unmarshal(contents, config); err != nil {
		return nil, err
	}

	return config, nil
}
