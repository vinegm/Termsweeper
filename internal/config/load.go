package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/vinegm/termsweeper/internal/storage/paths"
)

// validate reports an error for any setting
// the rest of the program can't work with.
func (config Config) validate() error {
	if config.Behaviour.Scrolloff < 0 {
		return fmt.Errorf("behaviour.scrolloff: must not be negative, got %d", config.Behaviour.Scrolloff)
	}

	if config.Behaviour.SideScrolloff < 0 {
		return fmt.Errorf("behaviour.sidescrolloff: must not be negative, got %d", config.Behaviour.SideScrolloff)
	}

	return nil
}

// Loads configuration from file. The returned config is always usable:
// on error the defaults are returned alongside it.
func Load(path string) (Config, error) {
	config := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Default(), err
	}

	if err := config.validate(); err != nil {
		return Default(), err
	}

	return config, nil
}

// DefaultPath is the config file's path.
func DefaultPath() string {
	return paths.Config("config.yaml")
}
