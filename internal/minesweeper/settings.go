package minesweeper

import "fmt"

// Limits a board has to stay within to be playable.
const (
	MinDim   = 5
	MaxDim   = 1000 // only used to prevent a user from asking for a board too big to allocate
	MinMines = 1
)

// Settings for the game board.
type Settings struct {
	Rows  int
	Cols  int
	Mines int
	Seed  string // mine layout seed
}

// MaxMines is the most mines a board this size can hold: the first revealed
// cell is never mined, so at least one cell has to stay free.
func (settings Settings) MaxMines() int {
	return settings.Rows*settings.Cols - 1
}

// Validate reports an error for settings no game can be played with.
func (settings Settings) Validate() error {
	if settings.Rows < MinDim || settings.Rows > MaxDim {
		return fmt.Errorf("rows: must be between %d and %d, got %d", MinDim, MaxDim, settings.Rows)
	}

	if settings.Cols < MinDim || settings.Cols > MaxDim {
		return fmt.Errorf("cols: must be between %d and %d, got %d", MinDim, MaxDim, settings.Cols)
	}

	if settings.Mines < MinMines || settings.Mines > settings.MaxMines() {
		return fmt.Errorf("mines: must be between %d and %d, got %d", MinMines, settings.MaxMines(), settings.Mines)
	}

	return nil
}
