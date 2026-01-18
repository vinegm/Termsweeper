package src

import (
	lg "github.com/charmbracelet/lipgloss"
)

const (
	bgColor         = ""
	evenSquareColor = "#1f2430"
	oddSquareColor  = "#2a2e3a"
	selectedColor   = "#f7768e"
	borderColor     = "#000000"
)

var (
	evenSquareStyle = lg.NewStyle().Background(lg.Color(evenSquareColor)).Align(lg.Center)
	oddSquareStyle  = lg.NewStyle().Background(lg.Color(oddSquareColor)).Align(lg.Center)

	selectedSquareStyle = lg.NewStyle().Background(lg.Color(selectedColor)).Align(lg.Center).Bold(true)
	selectedMenuStyle   = lg.NewStyle().Foreground(lg.Color(selectedColor)).Align(lg.Center).Bold(true)

	textStyle = lg.NewStyle().Align(lg.Center)

	windowStyle = lg.NewStyle().
			Background(lg.Color(bgColor)).
			Padding(0, 1).
			Border(lg.RoundedBorder()).
			BorderForeground(lg.Color(borderColor)).
			Align(lg.Center)

	titleStyle = lg.NewStyle().Bold(true).Align(lg.Center)
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
