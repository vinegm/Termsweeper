package theme

import (
	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/termsweeper/internal/config"
)

// Theme holds every style the windows render
// with, resolved from the configuration.
type Theme struct {
	config config.Config

	highlightLabelNumberStyle lg.Style
	labelNumberStyle          lg.Style

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

	WindowStyle lg.Style
	ArrowStyle  lg.Style
}

// NewTheme resolves a Theme from the config.
func NewTheme(cfg config.Config) *Theme {
	theme := &Theme{config: cfg}
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

func styleFromConfig(base lg.Style, configStyle config.Style) lg.Style {
	return withBg(withFg(base, configStyle.Fg), configStyle.Bg)
}

// Updates styles based on the current configuration.
func (theme *Theme) setStyles() {
	colors := theme.config.Theme.Resolve()
	squares := colors.Squares

	base := lg.NewStyle().Align(lg.Center)

	theme.TextStyle = withFg(base, colors.Fg)
	theme.TitleStyle = theme.TextStyle.Bold(true)
	theme.SelectedMenuTextStyle = withFg(theme.TitleStyle, colors.Selected)
	theme.selectedSquareStyle = withBg(theme.TitleStyle, colors.Selected)

	theme.labelNumberStyle = withFg(base, colors.Label.Fg)
	theme.highlightLabelNumberStyle = withFg(base, colors.Label.Highlight)

	theme.evenSquareStyle = styleFromConfig(base, squares.Even)
	theme.oddSquareStyle = styleFromConfig(base, squares.Odd)
	theme.flaggedSquareStyle = styleFromConfig(base, squares.Flagged)
	theme.minedSquareStyle = styleFromConfig(base, squares.Mined)
	theme.explodedSquareStyle = styleFromConfig(base, squares.Exploded)

	theme.squareMineHintColor = make([]lg.Style, len(squares.Hints))
	for i, hintConfig := range squares.Hints {
		theme.squareMineHintColor[i] = styleFromConfig(base, hintConfig)
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
