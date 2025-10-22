package config

import (
	"fmt"
	"os"
	
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Connectors []ConnectorConfig `yaml:"connectors"`
}

type ConnectorConfig struct {
	Type       string `yaml:"type"`
	Executable string `yaml:"executable,omitempty"`
	Path       string `yaml:"path"`
}

func LoadConfiguration(path string) (*Configuration, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	var config Configuration
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	
	return &config, nil
}

func (c *Configuration) Validate() error {
	if len(c.Connectors) == 0 {
		return fmt.Errorf("no connectors configured")
	}
	
	for i, conn := range c.Connectors {
		if conn.Type == "" {
			return fmt.Errorf("connector %d: type is required", i)
		}
		if conn.Path == "" {
			return fmt.Errorf("connector %d: path is required", i)
		}
	}
	
	return nil
}