package tui

import (
	"fmt"

	lg "github.com/charmbracelet/lipgloss"

	"termsweeper/internal/config"
	"termsweeper/internal/minesweeper"
)

// Holds every style the windows render with, resolved from the configuration.
type Theme struct {
	config config.Config

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
}

func NewTheme(cfg config.Config) *Theme {
	theme := &Theme{config: cfg}
	theme.setStyles()

	return theme
}

func styleFromConfig(baseStyle lg.Style, configStyle config.SquareStyleConfig) lg.Style {
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
func (theme *Theme) setStyles() {
	genericStyle := lg.NewStyle().Foreground(lg.Color(theme.config.FgColor)).Align(lg.Center)
	genericMenuTextStyle := lg.NewStyle().Inherit(genericStyle).Bold(true)

	theme.textStyle = lg.NewStyle().Inherit(genericStyle)
	theme.titleStyle = lg.NewStyle().Inherit(genericMenuTextStyle)
	theme.selectedMenuTextStyle = lg.NewStyle().Inherit(genericMenuTextStyle).Foreground(lg.Color(theme.config.SelectedColor))
	theme.selectedSquareStyle = lg.NewStyle().Inherit(genericMenuTextStyle).Background(lg.Color(theme.config.SelectedColor))

	theme.evenSquareStyle = styleFromConfig(genericStyle, theme.config.EvenSquareStyle)
	theme.oddSquareStyle = styleFromConfig(genericStyle, theme.config.OddSquareStyle)
	theme.flaggedSquareStyle = styleFromConfig(genericStyle, theme.config.FlaggedSquareStyle)
	theme.minedSquareStyle = styleFromConfig(genericStyle, theme.config.MinedSquareStyle)
	theme.explodedSquareStyle = styleFromConfig(theme.minedSquareStyle, theme.config.ExplodedSquareStyle)

	theme.squareMineHintColor = make([]lg.Style, len(theme.config.SquareMineHintStyle))
	for i, hintConfig := range theme.config.SquareMineHintStyle {
		theme.squareMineHintColor[i] = styleFromConfig(genericStyle, hintConfig)
	}

	theme.windowStyle = lg.NewStyle().
		Background(lg.Color(theme.config.BgColor)).
		Padding(0, 1).
		Border(lg.RoundedBorder()).
		BorderBackground(lg.Color(theme.config.BgColor)).
		BorderForeground(lg.Color(theme.config.BorderColor)).
		Align(lg.Center)
}

func (theme *Theme) validateGrid(style lg.Style, cell minesweeper.Cell, kind string) lg.Style {
	var preserve bool
	switch kind {
	case "revealed":
		preserve = theme.config.PreserveGridBg.Revealed

	case "flag":
		preserve = theme.config.PreserveGridBg.Flag

	case "mine":
		preserve = theme.config.PreserveGridBg.Mine

	default:
		preserve = false
	}

	if !preserve {
		return style
	}

	if cell.Is(minesweeper.Odd) {
		return style.Background(theme.oddSquareStyle.GetBackground())
	}

	return style.Background(theme.evenSquareStyle.GetBackground())

}

// Returns the style for a cell based on its position and selection state.
func (theme *Theme) cellStyle(cell minesweeper.Cell, isSelected bool) lg.Style {
	if isSelected {
		return theme.selectedSquareStyle
	}

	if cell.Is(minesweeper.Revealed) && cell.Is(minesweeper.Mined) {
		if cell.Is(minesweeper.Exploded) {
			return theme.validateGrid(theme.explodedSquareStyle, cell, "mine")
		}
		return theme.validateGrid(theme.minedSquareStyle, cell, "mine")
	}

	if cell.Is(minesweeper.Flagged) {
		return theme.validateGrid(theme.flaggedSquareStyle, cell, "flag")
	}

	if cell.Is(minesweeper.Revealed) {
		if cell.Adj >= len(theme.config.SquareMineHintStyle) {
			return theme.textStyle
		}

		return theme.validateGrid(theme.squareMineHintColor[cell.Adj], cell, "revealed")
	}

	if cell.Is(minesweeper.Odd) {
		return theme.oddSquareStyle
	}

	return theme.evenSquareStyle
}

// Returns the char for a given cell based on its state.
func (theme *Theme) cellChar(cell minesweeper.Cell) string {
	if cell.Is(minesweeper.Revealed) {
		if cell.Is(minesweeper.Mined) {
			return theme.config.MineChar
		}

		if cell.Adj == 0 {
			return " "
		}

		return fmt.Sprintf("%d", cell.Adj)
	}

	if cell.Is(minesweeper.Flagged) {
		return theme.config.FlagChar
	}

	return theme.config.UnrevealedChar
}
