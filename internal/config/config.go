package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// EdgeSetting picks which edge(s) of the board a feature applies to. Row
// number labels and row highlighting are horizontal, so they take "left" or
// "right"; column number labels and column highlighting are vertical, so they
// take "top" or "bottom". "both" covers either pair, and "" disables the
// feature.
type EdgeSetting string

const (
	EdgeNone   EdgeSetting = ""
	EdgeBoth   EdgeSetting = "both"
	EdgeLeft   EdgeSetting = "left"
	EdgeRight  EdgeSetting = "right"
	EdgeTop    EdgeSetting = "top"
	EdgeBottom EdgeSetting = "bottom"
)

func (e EdgeSetting) Left() bool   { return e == EdgeLeft || e == EdgeBoth }
func (e EdgeSetting) Right() bool  { return e == EdgeRight || e == EdgeBoth }
func (e EdgeSetting) Top() bool    { return e == EdgeTop || e == EdgeBoth }
func (e EdgeSetting) Bottom() bool { return e == EdgeBottom || e == EdgeBoth }

// Any reports whether the feature is enabled on at least one edge.
func (e EdgeSetting) Any() bool { return e != EdgeNone }

var (
	horizontalEdges = map[EdgeSetting]bool{EdgeNone: true, EdgeBoth: true, EdgeLeft: true, EdgeRight: true}
	verticalEdges   = map[EdgeSetting]bool{EdgeNone: true, EdgeBoth: true, EdgeTop: true, EdgeBottom: true}
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

	// Minimum rows/columns kept visible between the cursor and the viewport edge
	// when scrolling, like nvim's 'scrolloff'/'sidescrolloff'
	Scrolloff     int `yaml:"scrolloff"`
	SideScrolloff int `yaml:"sidescrolloff"`

	// Character settings
	FlagChar       string `yaml:"flag_char"`
	MineChar       string `yaml:"mine_char"`
	UnrevealedChar string `yaml:"unrevealed_char"`

	// Row/column coordinate labels around the board: which edge(s) show them
	LineNumbers   EdgeSetting `yaml:"line_numbers"`
	ColumnNumbers EdgeSetting `yaml:"column_numbers"`

	// Show line/column numbers relative to the cursor's position, like nvim's relativenumber
	RelativeLineNumbers   bool `yaml:"relative_line_numbers"`
	RelativeColumnNumbers bool `yaml:"relative_column_numbers"`

	// Colors for line/column number labels
	LabelNumberColor          string `yaml:"label_number_color"`           // labels not on the cursor's row/column
	HighlightLabelNumberColor string `yaml:"highlight_label_number_color"` // label on the cursor's row/column

	// Style settings
	FgColor       string `yaml:"fg_color"`
	BgColor       string `yaml:"bg_color"`
	SelectedColor string `yaml:"selected_color"`
	BorderColor   string `yaml:"border_color"`

	// Cursor row/column highlight, split at the cursor: which side(s) light up
	HighlightRow       EdgeSetting `yaml:"highlight_row"`
	HighlightCol       EdgeSetting `yaml:"highlight_col"`
	HighlightEvenColor string      `yaml:"highlight_even_color"`
	HighlightOddColor  string      `yaml:"highlight_odd_color"`

	// Board styles
	EvenSquareStyle     SquareStyleConfig   `yaml:"even_square_style"`
	OddSquareStyle      SquareStyleConfig   `yaml:"odd_square_style"`
	SquareMineHintStyle []SquareStyleConfig `yaml:"square_mine_hint_style"` // indexed by number of adjacent mines
	FlaggedSquareStyle  SquareStyleConfig   `yaml:"flagged_square_style"`
	MinedSquareStyle    SquareStyleConfig   `yaml:"mined_square_style"`
	ExplodedSquareStyle SquareStyleConfig   `yaml:"exploded_square_style"`

	PreserveGridBg PreserveGridBgConfig `yaml:"preserve_grid_bg"`
}

// The built-in configuration, used as the base the config file is layered over.
func Default() Config {
	return Config{
		BoardRows: 16,
		BoardCols: 16,
		MineCount: 40,

		Scrolloff:     0,
		SideScrolloff: 0,

		FlagChar:       "F",
		MineChar:       "X",
		UnrevealedChar: " ",

		LineNumbers:   EdgeNone,
		ColumnNumbers: EdgeNone,

		RelativeLineNumbers:   false,
		RelativeColumnNumbers: false,

		LabelNumberColor:          "",
		HighlightLabelNumberColor: "",

		FgColor:       "",
		SelectedColor: "#5a5a5a",
		BgColor:       "",
		BorderColor:   "",

		HighlightRow:       EdgeNone,
		HighlightCol:       EdgeNone,
		HighlightEvenColor: "#3a3a3a",
		HighlightOddColor:  "#454545",

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
}

// validate reports an error for any EdgeSetting that was given a value from
// the wrong axis, e.g. "top" for a left/right setting.
func (config Config) validate() error {
	for _, field := range []struct {
		name  string
		value EdgeSetting
		valid map[EdgeSetting]bool
	}{
		{"line_numbers", config.LineNumbers, horizontalEdges},
		{"highlight_row", config.HighlightRow, horizontalEdges},
		{"column_numbers", config.ColumnNumbers, verticalEdges},
		{"highlight_col", config.HighlightCol, verticalEdges},
	} {
		if !field.valid[field.value] {
			return fmt.Errorf("%s: invalid value %q", field.name, field.value)
		}
	}

	if config.Scrolloff < 0 {
		return fmt.Errorf("scrolloff: must not be negative, got %d", config.Scrolloff)
	}

	if config.SideScrolloff < 0 {
		return fmt.Errorf("sidescrolloff: must not be negative, got %d", config.SideScrolloff)
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

func DefaultPath() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config", "termsweeper", "config.yaml")
	}

	// crappy fallback, might not even work
	return "$HOME/.config/termsweeper/config.yaml"
}
