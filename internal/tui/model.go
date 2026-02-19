package tui

import (
	tea "charm.land/bubbletea/v2"

	"termsweeper/internal/config"
	"termsweeper/internal/minesweeper"
)

// Board settings the player can change from the menu. Seeded from the
// configuration, session-local from then on.
type Settings struct {
	Rows  int
	Cols  int
	Mines int
}

// Bubble Tea application state.
type Model struct {
	game          *minesweeper.Game // game state
	settings      Settings          // board settings
	theme         *Theme            // resolved styles
	width, height int               // terminal dimensions

	CurrentWindow Window
	MenuWin       *MenuWindow
	BoardWin      *BoardWindow
}

func (model Model) Init() tea.Cmd {
	return tick()
}

// Creates a new model with initial state.
func New(cfg config.Config) Model {
	menuWinInst := NewMenuWindow()

	model := Model{
		settings:      Settings{Rows: cfg.BoardRows, Cols: cfg.BoardCols, Mines: cfg.MineCount},
		theme:         NewTheme(cfg),
		CurrentWindow: menuWinInst,
		MenuWin:       menuWinInst,
		BoardWin:      NewBoardWindow(),
	}
	model.newGame()

	return model
}

// Replaces the board with a fresh one built from the current settings.
func (model *Model) newGame() {
	model.game = minesweeper.New(model.settings.Rows, model.settings.Cols, model.settings.Mines)
}
