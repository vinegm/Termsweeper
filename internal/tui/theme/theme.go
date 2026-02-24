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

// Theme holds every style the windows render with, resolved from the configuration.
type Theme struct {
	config config.Config

	LineNumbers           config.EdgeSetting
	ColumnNumbers         config.EdgeSetting
	RelativeLineNumbers   bool
	RelativeColumnNumbers bool
	HighlightRow          config.EdgeSetting
	HighlightCol          config.EdgeSetting

	Scrolloff     int
	SideScrolloff int

	TextStyle             lg.Style
	TitleStyle            lg.Style
	selectedSquareStyle   lg.Style
	SelectedMenuTextStyle lg.Style

	labelNumberStyle          lg.Style
	highlightLabelNumberStyle lg.Style

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
	theme := &Theme{
		config: cfg,

		LineNumbers:           cfg.LineNumbers,
		ColumnNumbers:         cfg.ColumnNumbers,
		RelativeLineNumbers:   cfg.RelativeLineNumbers,
		RelativeColumnNumbers: cfg.RelativeColumnNumbers,
		HighlightRow:          cfg.HighlightRow,
		HighlightCol:          cfg.HighlightCol,

		Scrolloff:     cfg.Scrolloff,
		SideScrolloff: cfg.SideScrolloff,
	}
	theme.setStyles()

	return theme
}

// Sets the foreground, leaving it inherited when no color is configured.
func withFg(base lg.Style, color string) lg.Style {
	if color == "" {
		return base
	}

	return base.Foreground(lg.Color(color))
}

// Sets the background, leaving it inherited when no color is configured.
func withBg(base lg.Style, color string) lg.Style {
	if color == "" {
		return base
	}

	return base.Background(lg.Color(color))
}

func styleFromConfig(baseStyle lg.Style, configStyle config.SquareStyleConfig) lg.Style {
	style := lg.NewStyle().Inherit(baseStyle)

	return withBg(withFg(style, configStyle.FgColor), configStyle.BgColor)
}

// Updates styles based on the current configuration.
func (theme *Theme) setStyles() {
	genericStyle := withFg(lg.NewStyle().Align(lg.Center), theme.config.FgColor)
	genericMenuTextStyle := lg.NewStyle().Inherit(genericStyle).Bold(true)

	theme.TextStyle = lg.NewStyle().Inherit(genericStyle)
	theme.TitleStyle = lg.NewStyle().Inherit(genericMenuTextStyle)
	theme.SelectedMenuTextStyle = withFg(lg.NewStyle().Inherit(genericMenuTextStyle), theme.config.SelectedColor)
	theme.selectedSquareStyle = withBg(lg.NewStyle().Inherit(genericMenuTextStyle), theme.config.SelectedColor)

	theme.labelNumberStyle = withFg(lg.NewStyle().Inherit(genericStyle), theme.config.LabelNumberColor)
	theme.highlightLabelNumberStyle = withFg(lg.NewStyle().Inherit(genericStyle), theme.config.HighlightLabelNumberColor)

	theme.evenSquareStyle = styleFromConfig(genericStyle, theme.config.EvenSquareStyle)
	theme.oddSquareStyle = styleFromConfig(genericStyle, theme.config.OddSquareStyle)
	theme.flaggedSquareStyle = styleFromConfig(genericStyle, theme.config.FlaggedSquareStyle)
	theme.minedSquareStyle = styleFromConfig(genericStyle, theme.config.MinedSquareStyle)
	theme.explodedSquareStyle = styleFromConfig(theme.minedSquareStyle, theme.config.ExplodedSquareStyle)

	theme.squareMineHintColor = make([]lg.Style, len(theme.config.SquareMineHintStyle))
	for i, hintConfig := range theme.config.SquareMineHintStyle {
		theme.squareMineHintColor[i] = styleFromConfig(genericStyle, hintConfig)
	}

	theme.ArrowStyle = withBg(withFg(lg.NewStyle(), theme.config.BorderColor), theme.config.BgColor)

	theme.WindowStyle = withBg(lg.NewStyle(), theme.config.BgColor).
		Padding(0, 1).
		Border(lg.RoundedBorder()).
		BorderBackground(lg.Color(theme.config.BgColor)).
		BorderForeground(lg.Color(theme.config.BorderColor)).
		Align(lg.Center)
}

// Restores the checkerboard background on a square whose kind is configured to keep it.
func (theme *Theme) preserveGridBg(style lg.Style, cell minesweeper.Cell, preserve bool) lg.Style {
	if !preserve {
		return style
	}

	if cell.StateIs(minesweeper.Odd) {
		return style.Background(theme.oddSquareStyle.GetBackground())
	}

	return style.Background(theme.evenSquareStyle.GetBackground())
}

// Returns the style for a revealed square, by its adjacent mine count.
func (theme *Theme) hintStyle(cell minesweeper.Cell) lg.Style {
	if cell.Adj < len(theme.config.SquareMineHintStyle) {
		return theme.squareMineHintColor[cell.Adj]
	}

	return theme.TextStyle
}

// Returns the background a highlighted square takes, keeping the checkerboard.
func (theme *Theme) highlightColor(cell minesweeper.Cell) string {
	if cell.StateIs(minesweeper.Odd) {
		return theme.config.HighlightOddColor
	}

	return theme.config.HighlightEvenColor
}

// Returns the style for a cell based on its position, selection and highlight state.
func (theme *Theme) CellStyle(cell minesweeper.Cell, isSelected, isHighlighted bool) lg.Style {
	if isSelected {
		return theme.selectedSquareStyle
	}

	var style lg.Style
	switch {
	case cell.StateIs(minesweeper.Exploded):
		style = theme.preserveGridBg(theme.explodedSquareStyle, cell, theme.config.PreserveGridBg.Mine)

	case cell.StateIs(minesweeper.Revealed) && cell.StateIs(minesweeper.Mined):
		style = theme.preserveGridBg(theme.minedSquareStyle, cell, theme.config.PreserveGridBg.Mine)

	case cell.StateIs(minesweeper.Flagged):
		style = theme.preserveGridBg(theme.flaggedSquareStyle, cell, theme.config.PreserveGridBg.Flag)

	case cell.StateIs(minesweeper.Revealed):
		style = theme.preserveGridBg(theme.hintStyle(cell), cell, theme.config.PreserveGridBg.Revealed)

	case cell.StateIs(minesweeper.Odd):
		style = theme.oddSquareStyle

	default:
		style = theme.evenSquareStyle
	}

	if isHighlighted {
		style = withBg(style, theme.highlightColor(cell))
	}

	return style
}

// LabelStyle returns the style for a row/column number label, based on whether
// it sits on the cursor's row/column.
func (theme *Theme) LabelStyle(idx, cursorIdx int) lg.Style {
	if idx == cursorIdx {
		return theme.highlightLabelNumberStyle
	}

	return theme.labelNumberStyle
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
