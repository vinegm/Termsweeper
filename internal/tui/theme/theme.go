package theme

import (
	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/termsweeper/internal/config"
)

// Theme holds every style the windows render
// with, resolved from the configuration.
type Theme struct {
	config config.Config

	Highlight                 Edge
	highlightLabelNumberStyle lg.Style

	Numbers               Edge
	labelNumberStyle      lg.Style
	RelativeLineNumbers   bool
	RelativeColumnNumbers bool

	TextStyle             lg.Style
	TitleStyle            lg.Style
	selectedSquareStyle   lg.Style
	SelectedMenuTextStyle lg.Style

	evenSquareStyle lg.Style
	oddSquareStyle  lg.Style

	flaggedSquareStyle  lg.Style
	minedSquareStyle    lg.Style
	squareMineHintColor []lg.Style
	explodedSquareStyle lg.Style

	WindowStyle   lg.Style
	ArrowStyle    lg.Style
	Scrolloff     int
	SideScrolloff int
}

func NewTheme(cfg config.Config) *Theme {
	theme := &Theme{
		config: cfg,

		Numbers:               cfg.Layout.Numbers,
		RelativeLineNumbers:   cfg.Layout.RelativeLineNumbers,
		RelativeColumnNumbers: cfg.Layout.RelativeColumnNumbers,
		Highlight:             cfg.Layout.Highlight,

		Scrolloff:     cfg.Behaviour.Scrolloff,
		SideScrolloff: cfg.Behaviour.SideScrolloff,
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

func styleFromConfig(baseStyle lg.Style, configStyle config.Style) lg.Style {
	style := lg.NewStyle().Inherit(baseStyle)

	return withBg(withFg(style, configStyle.Fg), configStyle.Bg)
}

// Updates styles based on the current configuration.
func (theme *Theme) setStyles() {
	colors, squares := theme.config.Theme, theme.config.Theme.Squares

	genericStyle := withFg(lg.NewStyle().Align(lg.Center), colors.Fg)
	genericMenuTextStyle := lg.NewStyle().Inherit(genericStyle).Bold(true)

	theme.TextStyle = lg.NewStyle().Inherit(genericStyle)
	theme.TitleStyle = lg.NewStyle().Inherit(genericMenuTextStyle)
	theme.SelectedMenuTextStyle = withFg(lg.NewStyle().Inherit(genericMenuTextStyle), colors.Selected)
	theme.selectedSquareStyle = withBg(lg.NewStyle().Inherit(genericMenuTextStyle), colors.Selected)

	theme.labelNumberStyle = withFg(lg.NewStyle().Inherit(genericStyle), colors.Label.Fg)
	theme.highlightLabelNumberStyle = withFg(lg.NewStyle().Inherit(genericStyle), colors.Label.Highlight)

	theme.evenSquareStyle = styleFromConfig(genericStyle, squares.Even)
	theme.oddSquareStyle = styleFromConfig(genericStyle, squares.Odd)
	theme.flaggedSquareStyle = styleFromConfig(genericStyle, squares.Flagged)
	theme.minedSquareStyle = styleFromConfig(genericStyle, squares.Mined)
	theme.explodedSquareStyle = styleFromConfig(theme.minedSquareStyle, squares.Exploded)

	theme.squareMineHintColor = make([]lg.Style, len(squares.Hints))
	for i, hintConfig := range squares.Hints {
		theme.squareMineHintColor[i] = styleFromConfig(genericStyle, hintConfig)
	}

	theme.ArrowStyle = withBg(withFg(lg.NewStyle(), colors.Border), colors.Bg)

	theme.WindowStyle = withBg(lg.NewStyle(), colors.Bg).
		Padding(0, 1).
		Border(lg.RoundedBorder()).
		BorderBackground(lg.Color(colors.Bg)).
		BorderForeground(lg.Color(colors.Border)).
		Align(lg.Center)
}

// LabelStyle returns the style for a row/column number label, based on whether
// it sits on the cursor's row/column.
func (theme *Theme) LabelStyle(idx, cursorIdx int) lg.Style {
	if idx == cursorIdx {
		return theme.highlightLabelNumberStyle
	}

	return theme.labelNumberStyle
}
