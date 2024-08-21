package conf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// This variable will store the configs as long as the server is alive
var cfg = new(Config)

// GetConfig is a getter for the private cfg variable
func GetConfig() *Config {
	return cfg
}

// Init initializes the config values by reading from the yaml file that is specified
func Init(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf(errorReadConfig, err)
	}

	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return fmt.Errorf(errYamlUnmarshal, err.Error())
	}
	return nil
}
