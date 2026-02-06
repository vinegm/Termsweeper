package src

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Square foreground and background color for a square.
type SquareStyleConfig struct {
	FgColor string `yaml:"fg_color"`
	BgColor string `yaml:"bg_color"`
}

type Config struct {
	// Game settings
	BoardRows int `yaml:"board_rows"`
	BoardCols int `yaml:"board_cols"`
	MineCount int `yaml:"mine_count"`

	// Character settings
	FlagChar string `yaml:"flag_char"`
	MineChar string `yaml:"mine_char"`

	// Style settings
	TextColor     string `yaml:"text_color"`
	BgColor       string `yaml:"bg_color"`
	SelectedColor string `yaml:"selected_color"`
	BorderColor   string `yaml:"border_color"`

	// Board styles
	EvenSquareStyle     SquareStyleConfig   `yaml:"even_square_style"`
	OddSquareStyle      SquareStyleConfig   `yaml:"odd_square_style"`
	SquareMineHintStyle []SquareStyleConfig `yaml:"square_mine_hint_style"` // indexed by number of adjacent mines
	FlaggedSquareStyle  SquareStyleConfig   `yaml:"flagged_square_style"`
	MinedSquareStyle    SquareStyleConfig   `yaml:"mined_square_style"`
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
		BoardRows: 16,
		BoardCols: 16,
		MineCount: 40,

		FlagChar: "F",
		MineChar: "X",

		TextColor:     "",
		SelectedColor: "#5a5a5a",
		BgColor:       "",
		BorderColor:   "",

		EvenSquareStyle: SquareStyleConfig{
			FgColor: "",
			BgColor: "#1C2B1C",
		},

		OddSquareStyle: SquareStyleConfig{
			FgColor: "",
			BgColor: "#263726",
		},

		FlaggedSquareStyle: SquareStyleConfig{
			FgColor: "",
			BgColor: "#FFD580",
		},

		MinedSquareStyle: SquareStyleConfig{
			FgColor: "",
			BgColor: "#FF0000",
		},

		SquareMineHintStyle: []SquareStyleConfig{
			{FgColor: "", BgColor: ""}, // 0 mines
			{FgColor: "", BgColor: ""}, // 1 mines
			{FgColor: "", BgColor: ""}, // 2 mines
			{FgColor: "", BgColor: ""}, // 3 mines
			{FgColor: "", BgColor: ""}, // 4 mines
			{FgColor: "", BgColor: ""}, // 5 mines
			{FgColor: "", BgColor: ""}, // 6 mines
			{FgColor: "", BgColor: ""}, // 7 mines
			{FgColor: "", BgColor: ""}, // 8 mines
		},
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
