package config

// Theme is every colour the windows render with.
type Theme struct {
	Fg       string `yaml:"fg"`       // default fg for all text and squares
	Bg       string `yaml:"bg"`       // default bg for the window and squares
	Selected string `yaml:"selected"` // fg of selected text, bg of the selected square
	Border   string `yaml:"border"`

	Label     Label   `yaml:"label"`
	Highlight Checker `yaml:"highlight"`
	Squares   Squares `yaml:"squares"`

	PreserveGridBg PreserveGridBg `yaml:"preserve_grid_bg"` // if we need to preserve the checkerboard background
}

// Style is a foreground and background colour.
type Style struct {
	Fg string `yaml:"fg"`
	Bg string `yaml:"bg"`
}

// Label is the row/column coordinate labels' colours.
type Label struct {
	Fg        string `yaml:"fg"`        // labels not on the cursor's row/column
	Highlight string `yaml:"highlight"` // label on the cursor's row/column
}

// Checker is a colour for each square of the checkerboard.
type Checker struct {
	Even string `yaml:"even"`
	Odd  string `yaml:"odd"`
}

// Squares is the board's colours, by square state.
type Squares struct {
	Even     Style   `yaml:"even"`
	Odd      Style   `yaml:"odd"`
	Flagged  Style   `yaml:"flagged"`
	Mined    Style   `yaml:"mined"`
	Exploded Style   `yaml:"exploded"`
	Hints    []Style `yaml:"hints"` // indexed by number of adjacent mines
}

// PreserveGridBg determines whether the checkerboard
// background is preserved for a given square state.
type PreserveGridBg struct {
	Revealed bool `yaml:"revealed"`
	Flag     bool `yaml:"flag"`
	Mine     bool `yaml:"mine"`
}
