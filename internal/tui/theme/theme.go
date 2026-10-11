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

		Numbers:               cfg.Numbers,
		RelativeLineNumbers:   cfg.RelativeLineNumbers,
		RelativeColumnNumbers: cfg.RelativeColumnNumbers,
		Highlight:             cfg.Highlight,

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

// LabelStyle returns the style for a row/column number label, based on whether
// it sits on the cursor's row/column.
func (theme *Theme) LabelStyle(idx, cursorIdx int) lg.Style {
	if idx == cursorIdx {
		return theme.highlightLabelNumberStyle
	}

	return theme.labelNumberStyle
}
