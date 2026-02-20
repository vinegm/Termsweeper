package tui

import (
	tea "charm.land/bubbletea/v2"

	"termsweeper/internal/config"
	"termsweeper/internal/minesweeper"
)

// Model represents the state of the application.
type Model struct {
	theme  *Theme
	screen *Screen

	currentWindow Window
	menu          *MenuWindow
	board         *BoardWindow
}

// New creates a new model with initial state.
func New(cfg config.Config) Model {
	theme := NewTheme(cfg)
	screen := &Screen{}
	gameSettings := minesweeper.Settings{
		Rows:  cfg.BoardRows,
		Cols:  cfg.BoardCols,
		Mines: cfg.MineCount,
	}

	menu := NewMenuWindow(theme, gameSettings)

	return Model{
		theme:         theme,
		screen:        screen,
		currentWindow: menu,
		menu:          menu,
		board:         NewBoardWindow(theme, screen),
	}
}

// Init initializes the model and returns an initial command.
//
//	WARN: This function is required by the Bubble Tea framework.
func (model Model) Init() tea.Cmd {
	return nil
}
