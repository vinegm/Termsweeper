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

type PreserveGridBgConfig struct {
	Revealed bool `yaml:"revealed"`
	Flag     bool `yaml:"flag"`
	Mine     bool `yaml:"mine"`
}

type Config struct {
	// Game settings
	BoardRows int `yaml:"board_rows"`
	BoardCols int `yaml:"board_cols"`
	MineCount int `yaml:"mine_count"`

	// Character settings
	FlagChar       string `yaml:"flag_char"`
	MineChar       string `yaml:"mine_char"`
	UnrevealedChar string `yaml:"unrevealed_char"`

	// Style settings
	FgColor       string `yaml:"fg_color"`
	BgColor       string `yaml:"bg_color"`
	SelectedColor string `yaml:"selected_color"`
	BorderColor   string `yaml:"border_color"`

	// Board styles
	EvenSquareStyle     SquareStyleConfig   `yaml:"even_square_style"`
	OddSquareStyle      SquareStyleConfig   `yaml:"odd_square_style"`
	SquareMineHintStyle []SquareStyleConfig `yaml:"square_mine_hint_style"` // indexed by number of adjacent mines
	FlaggedSquareStyle  SquareStyleConfig   `yaml:"flagged_square_style"`
	MinedSquareStyle    SquareStyleConfig   `yaml:"mined_square_style"`
	ExplodedSquareStyle SquareStyleConfig   `yaml:"exploded_square_style"`

	PreserveGridBg PreserveGridBgConfig `yaml:"preserve_grid_bg"`
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

		FlagChar:       "F",
		MineChar:       "X",
		UnrevealedChar: " ",

		FgColor:       "",
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
			FgColor: "#FFD580",
			BgColor: "",
		},

		MinedSquareStyle: SquareStyleConfig{
			FgColor: "#800000",
			BgColor: "",
		},

		ExplodedSquareStyle: SquareStyleConfig{
			FgColor: "",
			BgColor: "",
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

		PreserveGridBg: PreserveGridBgConfig{
			Revealed: false,
			Flag:     true,
			Mine:     true,
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
