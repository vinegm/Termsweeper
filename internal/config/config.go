package config

// Square foreground and background color for a square.
type SquareStyleConfig struct {
	FgColor string `yaml:"fg_color"`
	BgColor string `yaml:"bg_color"`
}

// PreserveGridBgConfig determines whether the checkerboard
// background is preserved for a given square state.
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

	// Directions out from the cursor its row/column highlight extends in
	Highlight                 EdgeSetting `yaml:"highlight"`
	HighlightLabelNumberColor string      `yaml:"highlight_label_number_color"` // label on the cursor's row/column
	HighlightEvenColor        string      `yaml:"highlight_even_color"`
	HighlightOddColor         string      `yaml:"highlight_odd_color"`

	// Board edges the row/column coordinate labels are drawn on
	Numbers               EdgeSetting `yaml:"numbers"`
	LabelNumberColor      string      `yaml:"label_number_color"`      // labels not on the cursor's row/column
	RelativeLineNumbers   bool        `yaml:"relative_line_numbers"`   // like nvim's 'relativenumber'
	RelativeColumnNumbers bool        `yaml:"relative_column_numbers"` // like nvim's 'relativenumber' but for columns

	// Style settings
	FgColor       string `yaml:"fg_color"`
	BgColor       string `yaml:"bg_color"`
	SelectedColor string `yaml:"selected_color"`
	BorderColor   string `yaml:"border_color"`

	// Board styles
	EvenSquareStyle SquareStyleConfig `yaml:"even_square_style"`
	OddSquareStyle  SquareStyleConfig `yaml:"odd_square_style"`

	FlaggedSquareStyle  SquareStyleConfig   `yaml:"flagged_square_style"`
	MinedSquareStyle    SquareStyleConfig   `yaml:"mined_square_style"`
	SquareMineHintStyle []SquareStyleConfig `yaml:"square_mine_hint_style"` // indexed by number of adjacent mines
	ExplodedSquareStyle SquareStyleConfig   `yaml:"exploded_square_style"`

	PreserveGridBg PreserveGridBgConfig `yaml:"preserve_grid_bg"` // if we need to preserve the checkerboard background

	// Character settings
	FlagChar       string `yaml:"flag_char"`
	MineChar       string `yaml:"mine_char"`
	UnrevealedChar string `yaml:"unrevealed_char"`

	// Minimum rows/columns kept visible between the
	// cursor and the viewport edge when scrolling
	Scrolloff     int `yaml:"scrolloff"`     // like nvim's 'scrolloff'
	SideScrolloff int `yaml:"sidescrolloff"` // like nvim's 'sidescrolloff'
}

// The built-in configuration, used as the
// base the config file is layered over.
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

		Numbers: EdgeNone,

		RelativeLineNumbers:   false,
		RelativeColumnNumbers: false,

		LabelNumberColor:          "",
		HighlightLabelNumberColor: "",

		FgColor:       "",
		SelectedColor: "#5a5a5a",
		BgColor:       "",
		BorderColor:   "",

		Highlight:          EdgeNone,
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
