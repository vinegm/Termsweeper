package theme

import (
	"fmt"

	lg "charm.land/lipgloss/v2"

	"termsweeper/internal/config"
	"termsweeper/internal/minesweeper"
)

// Glyphs marking a window edge the board continues past.
const (
	ArrowUp    = "▲"
	ArrowDown  = "▼"
	ArrowLeft  = "◀"
	ArrowRight = "▶"
)

// Holds every style the windows render with, resolved from the configuration.
type Theme struct {
	config config.Config

	TextStyle             lg.Style
	TitleStyle            lg.Style
	selectedSquareStyle   lg.Style
	SelectedMenuTextStyle lg.Style

	evenSquareStyle     lg.Style
	oddSquareStyle      lg.Style
	flaggedSquareStyle  lg.Style
	minedSquareStyle    lg.Style
	explodedSquareStyle lg.Style

	squareMineHintColor []lg.Style

	ArrowStyle  lg.Style
	WindowStyle lg.Style
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

	theme.TextStyle = lg.NewStyle().Inherit(genericStyle)
	theme.TitleStyle = lg.NewStyle().Inherit(genericMenuTextStyle)
	theme.SelectedMenuTextStyle = lg.NewStyle().Inherit(genericMenuTextStyle).Foreground(lg.Color(theme.config.SelectedColor))
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

	theme.ArrowStyle = lg.NewStyle().
		Foreground(lg.Color(theme.config.BorderColor)).
		Background(lg.Color(theme.config.BgColor))

	theme.WindowStyle = lg.NewStyle().
		Background(lg.Color(theme.config.BgColor)).
		Padding(0, 1).
		Border(lg.RoundedBorder()).
		BorderBackground(lg.Color(theme.config.BgColor)).
		BorderForeground(lg.Color(theme.config.BorderColor)).
		Align(lg.Center)
}

// The square kinds that can be configured to keep the checkerboard background.
type squareKind int

const (
	revealedSquare squareKind = iota
	flagSquare
	mineSquare
)

// Restores the checkerboard background on a square whose kind is configured to keep it.
func (theme *Theme) preserveGridBg(style lg.Style, cell minesweeper.Cell, kind squareKind) lg.Style {
	var preserve bool
	switch kind {
	case revealedSquare:
		preserve = theme.config.PreserveGridBg.Revealed

	case flagSquare:
		preserve = theme.config.PreserveGridBg.Flag

	case mineSquare:
		preserve = theme.config.PreserveGridBg.Mine
	}

	if !preserve {
		return style
	}

	if cell.StateIs(minesweeper.Odd) {
		return style.Background(theme.oddSquareStyle.GetBackground())
	}

	return style.Background(theme.evenSquareStyle.GetBackground())
}

// Returns the style for a cell based on its position and selection state.
func (theme *Theme) CellStyle(cell minesweeper.Cell, isSelected bool) lg.Style {
	if isSelected {
		return theme.selectedSquareStyle
	}

	if cell.StateIs(minesweeper.Revealed) && cell.StateIs(minesweeper.Mined) {
		if cell.StateIs(minesweeper.Exploded) {
			return theme.preserveGridBg(theme.explodedSquareStyle, cell, mineSquare)
		}
		return theme.preserveGridBg(theme.minedSquareStyle, cell, mineSquare)
	}

	if cell.StateIs(minesweeper.Flagged) {
		return theme.preserveGridBg(theme.flaggedSquareStyle, cell, flagSquare)
	}

	if cell.StateIs(minesweeper.Revealed) {
		if cell.Adj >= len(theme.config.SquareMineHintStyle) {
			return theme.TextStyle
		}

		return theme.preserveGridBg(theme.squareMineHintColor[cell.Adj], cell, revealedSquare)
	}

	if cell.StateIs(minesweeper.Odd) {
		return theme.oddSquareStyle
	}

	return theme.evenSquareStyle
}

// Returns the char for a given cell based on its state.
func (theme *Theme) CellChar(cell minesweeper.Cell) string {
	if cell.StateIs(minesweeper.Revealed) {
		if cell.StateIs(minesweeper.Mined) {
			return theme.config.MineChar
		}

		if cell.Adj == 0 {
			return " "
		}

		return fmt.Sprintf("%d", cell.Adj)
	}

	if cell.StateIs(minesweeper.Flagged) {
		return theme.config.FlagChar
	}

	return theme.config.UnrevealedChar
}
