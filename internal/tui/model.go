package tui

import (
	tea "charm.land/bubbletea/v2"

	"termsweeper/internal/config"
	"termsweeper/internal/minesweeper"
	"termsweeper/internal/tui/board"
	"termsweeper/internal/tui/menu"
	"termsweeper/internal/tui/theme"
)

// Model represents the state of the application.
type Model struct {
	width  int
	height int

	currentWindow Window
	menu          *menu.Window
	board         *board.Window
}

// New creates a new model with initial state.
func New(cfg config.Config) Model {
	th := theme.NewTheme(cfg)
	gameSettings := minesweeper.Settings{
		Rows:  cfg.BoardRows,
		Cols:  cfg.BoardCols,
		Mines: cfg.MineCount,
	}

	menuWindow := menu.New(th, gameSettings)

	return Model{
		currentWindow: menuWindow,
		menu:          menuWindow,
		board:         board.New(th),
	}
}

// Init initializes the model and returns an initial command.
//
//	WARN: This function is required by the Bubble Tea framework.
func (model Model) Init() tea.Cmd {
	return nil
}
