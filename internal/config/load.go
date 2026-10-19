package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/vinegm/termsweeper/internal/storage/paths"
)

// clamp brings each setting the rest of the program can't work with
// back into range, keeping the rest of the file.
func (config *Config) clamp() {
	config.Behaviour.Scrolloff = max(config.Behaviour.Scrolloff, 0)
	config.Behaviour.SideScrolloff = max(config.Behaviour.SideScrolloff, 0)
}

// Load reads the config file at path. The returned config is always
// usable: on error the defaults are returned alongside it.
func Load(path string) (Config, error) {
	config := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Default(), err
	}

	config.clamp()

	return config, nil
}

// DefaultPath is the config file's path.
func DefaultPath() string {
	return paths.Config("config.yaml")
}
