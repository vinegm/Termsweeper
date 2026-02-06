package src

import (
	lg "github.com/charmbracelet/lipgloss"
)

var (
	textStyle             lg.Style
	titleStyle            lg.Style
	selectedSquareStyle   lg.Style
	selectedMenuTextStyle lg.Style

	evenSquareStyle    lg.Style
	oddSquareStyle     lg.Style
	flaggedSquareStyle lg.Style
	minedSquareStyle   lg.Style

	squareMineHintColor []lg.Style

	windowStyle lg.Style
)

func styleFromConfig(baseStyle lg.Style, configStyle SquareStyleConfig) lg.Style {
	style := lg.NewStyle().Inherit(baseStyle)
	if configStyle.FgColor != "" {
		style = style.Foreground(lg.Color(configStyle.FgColor))
	}

	if configStyle.BgColor != "" {
		style = style.Background(lg.Color(configStyle.BgColor))
	}

	return style
}

// Updates styles based on the current configuration.
func setStyles() {
	genericStyle := lg.NewStyle().Foreground(lg.Color(AppConfig.TextColor)).Align(lg.Center)
	genericMenuTextStyle := lg.NewStyle().Inherit(genericStyle).Bold(true)

	textStyle = lg.NewStyle().Inherit(genericStyle)
	titleStyle = lg.NewStyle().Inherit(genericMenuTextStyle)
	selectedMenuTextStyle = lg.NewStyle().Inherit(genericMenuTextStyle).Foreground(lg.Color(AppConfig.SelectedColor))
	selectedSquareStyle = lg.NewStyle().Inherit(genericMenuTextStyle).Background(lg.Color(AppConfig.SelectedColor))

	evenSquareStyle = styleFromConfig(genericStyle, AppConfig.EvenSquareStyle)
	oddSquareStyle = styleFromConfig(genericStyle, AppConfig.OddSquareStyle)
	flaggedSquareStyle = styleFromConfig(genericStyle, AppConfig.FlaggedSquareStyle)
	minedSquareStyle = styleFromConfig(genericStyle, AppConfig.MinedSquareStyle)

	squareMineHintColor = make([]lg.Style, len(AppConfig.SquareMineHintStyle))
	for i, hintConfig := range AppConfig.SquareMineHintStyle {
		squareMineHintColor[i] = styleFromConfig(genericStyle, hintConfig)
	}

	windowStyle = lg.NewStyle().
		Background(lg.Color(AppConfig.BgColor)).
		Padding(0, 1).
		Border(lg.RoundedBorder()).
		BorderBackground(lg.Color(AppConfig.BgColor)).
		BorderForeground(lg.Color(AppConfig.BorderColor)).
		Align(lg.Center)
}

// Returns the style for a cell based on its position and selection state.
func cellStyle(cell cell, isSelected bool) lg.Style {
	if isSelected {
		return selectedSquareStyle
	}

	if cell.is(flagged) {
		return flaggedSquareStyle
	}

	if cell.is(revealed) {
		if cell.is(mined) {
			return minedSquareStyle
		}

		if cell.adj >= len(squareMineHintColor) {
			return textStyle
		}

		return squareMineHintColor[cell.adj]
	}

	if cell.is(odd) {
		return oddSquareStyle
	}

	return evenSquareStyle
}
