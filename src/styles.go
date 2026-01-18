package src

import (
	lg "github.com/charmbracelet/lipgloss"
)

const (
	selectedColor = "#ff6b6b"
)

var (
	evenSquareStyle     = lg.NewStyle().Background(lg.Color("#2d2d2d"))
	oddSquareStyle      = lg.NewStyle().Background(lg.Color("#3a3a3a"))
	selectedSquareStyle = lg.NewStyle().Background(lg.Color(selectedColor))

	selectedMenuStyle   = lg.NewStyle().Foreground(lg.Color(selectedColor))
)

// Returns the style for a cell based on its position and selection state.
func cellStyle(row, col int, isSelected bool) lg.Style {
	if isSelected {
		return selectedSquareStyle
	}

	if (row+col)%2 == 0 {
		return evenSquareStyle
	}

	return oddSquareStyle
}
