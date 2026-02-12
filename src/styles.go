package src

import (
	lg "github.com/charmbracelet/lipgloss"
)

var (
	genericStyle lg.Style

	textStyle             lg.Style
	titleStyle            lg.Style
	selectedSquareStyle   lg.Style
	selectedMenuTextStyle lg.Style

	evenSquareStyle     lg.Style
	oddSquareStyle      lg.Style
	flaggedSquareStyle  lg.Style
	minedSquareStyle    lg.Style
	explodedSquareStyle lg.Style

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
	genericStyle = lg.NewStyle().Foreground(lg.Color(AppConfig.FgColor)).Align(lg.Center)
	genericMenuTextStyle := lg.NewStyle().Inherit(genericStyle).Bold(true)

	textStyle = lg.NewStyle().Inherit(genericStyle)
	titleStyle = lg.NewStyle().Inherit(genericMenuTextStyle)
	selectedMenuTextStyle = lg.NewStyle().Inherit(genericMenuTextStyle).Foreground(lg.Color(AppConfig.SelectedColor))
	selectedSquareStyle = lg.NewStyle().Inherit(genericMenuTextStyle).Background(lg.Color(AppConfig.SelectedColor))

	evenSquareStyle = styleFromConfig(genericStyle, AppConfig.EvenSquareStyle)
	oddSquareStyle = styleFromConfig(genericStyle, AppConfig.OddSquareStyle)
	flaggedSquareStyle = styleFromConfig(genericStyle, AppConfig.FlaggedSquareStyle)
	minedSquareStyle = styleFromConfig(genericStyle, AppConfig.MinedSquareStyle)
	explodedSquareStyle = styleFromConfig(minedSquareStyle, AppConfig.ExplodedSquareStyle)

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

func validateGrid(style lg.Style, cell cell, kind string) lg.Style {
	var preserve bool
	switch kind {
	case "revealed":
		preserve = AppConfig.PreserveGridBg.Revealed

	case "flag":
		preserve = AppConfig.PreserveGridBg.Flag

	case "mine":
		preserve = AppConfig.PreserveGridBg.Mine

	default:
		preserve = false
	}

	if !preserve {
		return style
	}

	if cell.is(odd) {
		return style.Background(oddSquareStyle.GetBackground())
	}

	return style.Background(evenSquareStyle.GetBackground())

}

// Returns the style for a cell based on its position and selection state.
func cellStyle(cell cell, isSelected bool) lg.Style {
	if isSelected {
		return selectedSquareStyle
	}

	if cell.is(revealed) && cell.is(mined) {
		if cell.is(exploded) {
			return validateGrid(explodedSquareStyle, cell, "mine")
		}
		return validateGrid(minedSquareStyle, cell, "mine")
	}

	if cell.is(flagged) {
		return validateGrid(flaggedSquareStyle, cell, "flag")
	}

	if cell.is(revealed) {
		if cell.adj >= len(AppConfig.SquareMineHintStyle) {
			return textStyle
		}

		return validateGrid(squareMineHintColor[cell.adj], cell, "revealed")
	}

	if cell.is(odd) {
		return oddSquareStyle
	}

	return evenSquareStyle
}
