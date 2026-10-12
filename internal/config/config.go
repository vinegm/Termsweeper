package config

// Config is config.yaml.
// Keys are documented in Themes/example-config.yaml.
type Config struct {
	Board     Board     `yaml:"board"`
	Behaviour Behaviour `yaml:"behaviour"`
	Layout    Layout    `yaml:"layout"`
	Chars     Chars     `yaml:"chars"`
	Theme     Theme     `yaml:"theme"`
}

// Board is the size of a new game.
type Board struct {
	Rows  int `yaml:"rows"`
	Cols  int `yaml:"cols"`
	Mines int `yaml:"mines"`
}

// Behaviour is how the cursor moves, not how it looks.
type Behaviour struct {
	// Minimum rows/columns kept visible between the
	// cursor and the viewport edge when scrolling
	Scrolloff     int `yaml:"scrolloff"`     // like nvim's 'scrolloff'
	SideScrolloff int `yaml:"sidescrolloff"` // like nvim's 'sidescrolloff'
}

// Layout is what is drawn around the cursor and the board.
type Layout struct {
	// Directions out from the cursor its row/column highlight extends in
	Highlight EdgeSetting `yaml:"highlight"`

	// Board edges the row/column coordinate labels are drawn on
	Numbers               EdgeSetting `yaml:"numbers"`
	RelativeLineNumbers   bool        `yaml:"relative_line_numbers"`   // like nvim's 'relativenumber'
	RelativeColumnNumbers bool        `yaml:"relative_column_numbers"` // like nvim's 'relativenumber' but for columns
}

// Chars is what a square shows in each state.
type Chars struct {
	Flag       string `yaml:"flag"`
	Mine       string `yaml:"mine"`
	Unrevealed string `yaml:"unrevealed"`
}

// The built-in configuration, used as the
// base the config file is layered over.
func Default() Config {
	return Config{
		Board: Board{
			Rows:  16,
			Cols:  16,
			Mines: 40,
		},

		Behaviour: Behaviour{
			Scrolloff:     0,
			SideScrolloff: 0,
		},

		Layout: Layout{
			Highlight:             EdgeNone,
			Numbers:               EdgeNone,
			RelativeLineNumbers:   false,
			RelativeColumnNumbers: false,
		},

		Chars: Chars{
			Flag:       "F",
			Mine:       "X",
			Unrevealed: " ",
		},

		Theme: Theme{
			Fg:       "",
			Bg:       "",
			Selected: "#5a5a5a",
			Border:   "",

			Label: Label{
				Fg:        "",
				Highlight: "",
			},

			Highlight: Checker{
				Even: "#3a3a3a",
				Odd:  "#454545",
			},

			Squares: Squares{
				Even:     Style{Fg: "", Bg: "#1C2B1C"},
				Odd:      Style{Fg: "", Bg: "#263726"},
				Flagged:  Style{Fg: "#FFD580", Bg: ""},
				Mined:    Style{Fg: "#800000", Bg: ""},
				Exploded: Style{Fg: "", Bg: ""},

				Hints: []Style{
					{Fg: "", Bg: ""}, // 0 mines
					{Fg: "", Bg: ""}, // 1 mines
					{Fg: "", Bg: ""}, // 2 mines
					{Fg: "", Bg: ""}, // 3 mines
					{Fg: "", Bg: ""}, // 4 mines
					{Fg: "", Bg: ""}, // 5 mines
					{Fg: "", Bg: ""}, // 6 mines
					{Fg: "", Bg: ""}, // 7 mines
					{Fg: "", Bg: ""}, // 8 mines
				},
			},

			PreserveGridBg: PreserveGridBg{
				Revealed: false,
				Flag:     true,
				Mine:     true,
			},
		},
	}
}
