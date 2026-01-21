package src

import (
	lg "github.com/charmbracelet/lipgloss"
)

var (
	evenSquareStyle lg.Style
	oddSquareStyle  lg.Style

	selectedSquareStyle lg.Style
	selectedMenuStyle   lg.Style

	textStyle   lg.Style
	windowStyle lg.Style
	titleStyle  lg.Style
)

// Updates styles based on the current configuration.
func setStyles() {
	evenSquareStyle = lg.NewStyle().Foreground(lg.Color(AppConfig.TextColor)).Background(lg.Color(AppConfig.EvenSquareColor)).Align(lg.Center)
	oddSquareStyle = lg.NewStyle().Foreground(lg.Color(AppConfig.TextColor)).Background(lg.Color(AppConfig.OddSquareColor)).Align(lg.Center)

	selectedSquareStyle = lg.NewStyle().Foreground(lg.Color(AppConfig.TextColor)).Background(lg.Color(AppConfig.SelectedColor)).Align(lg.Center).Bold(true)
	selectedMenuStyle = lg.NewStyle().Foreground(lg.Color(AppConfig.SelectedColor)).Align(lg.Center).Bold(true)
	textStyle = lg.NewStyle().Foreground(lg.Color(AppConfig.TextColor)).Align(lg.Center)

	windowStyle = lg.NewStyle().
		Background(lg.Color(AppConfig.BgColor)).
		Padding(0, 1).
		Border(lg.RoundedBorder()).
		BorderBackground(lg.Color(AppConfig.BgColor)).
		BorderForeground(lg.Color(AppConfig.BorderColor)).
		Align(lg.Center)

	titleStyle = lg.NewStyle().Foreground(lg.Color(AppConfig.TextColor)).Bold(true).Align(lg.Center)
}

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
