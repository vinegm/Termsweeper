package tui

import (
	tea "charm.land/bubbletea/v2"

	"termsweeper/internal/config"
	"termsweeper/internal/minesweeper"
)

// Model represents the state of the application.
type Model struct {
	theme         *Theme
	width, height int

	currentWindow Window
	menu          *MenuWindow
	board         *BoardWindow
}

// New creates a new model with initial state.
func New(cfg config.Config) Model {
	theme := NewTheme(cfg)
	gameSettings := minesweeper.Settings{
		Rows:  cfg.BoardRows,
		Cols:  cfg.BoardCols,
		Mines: cfg.MineCount,
	}

	menu := NewMenuWindow(theme, gameSettings)

	return Model{
		theme:         theme,
		currentWindow: menu,
		menu:          menu,
		board:         NewBoardWindow(theme),
	}
}

// Init initializes the model and returns an initial command.
//
//	WARN: This function is required by the Bubble Tea framework.
func (model Model) Init() tea.Cmd {
	return nil
}

// switchWindow switches the current window to the recieved one.
func (model Model) switchWindow(window Window) (Model, tea.Cmd) {
	var cmd tea.Cmd
	model.currentWindow, cmd = window.Update(
		tea.WindowSizeMsg{
			Width:  model.width,
			Height: model.height,
		},
	)

	return model, cmd
}
