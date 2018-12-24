package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	APIURL   string          `yaml:"apiurl"`
	Accounts []string        `yaml:"accounts,omitempty"`
	Tokens   []TokenContract `yaml:"tokens,omitempty"`
}

type TokenContract struct {
	Account string `yaml:"account"`
	Symbol  string `yaml:"symbol"`
}

func LoadConfig(confFile string) (*Config, error) {
	var c = &Config{}

	yamlFile, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading config file: %s", err)
	}

	if err := yaml.UnmarshalStrict(yamlFile, &c); err != nil {
		return nil, fmt.Errorf("Error parsing config file: %s", err)
	}

	return c, nil
}
