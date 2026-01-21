package src

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TextColor string `yaml:"text_color"`

	BgColor         string `yaml:"bg_color"`
	EvenSquareColor string `yaml:"even_square_color"`
	OddSquareColor  string `yaml:"odd_square_color"`
	SelectedColor   string `yaml:"selected_color"`
	BorderColor     string `yaml:"border_color"`

	FlagChar string `yaml:"flag_char"`
	MineChar string `yaml:"mine_char"`
}

var AppConfig Config

func configDefaultPath() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config", "termsweeper", "config.yaml")
	}

	// crappy fallback, might not even work
	return "$HOME/.config/termsweeper/config.yaml"
}

// Loads configuration from file.
func LoadConfig() error {
	AppConfig = Config{
		TextColor: "",

		BgColor:         "",
		EvenSquareColor: "#1C2B1C",
		OddSquareColor:  "#263726",
		SelectedColor:   "#FF0000",
		BorderColor:     "",

		FlagChar: "F",
		MineChar: "X",
	}

	path := configDefaultPath()
	data, err := os.ReadFile(path)
	if err != nil {
		setStyles()

		return err
	}

	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		return err
	}

	setStyles()

	return nil
}
