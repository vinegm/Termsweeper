package src

import tea "github.com/charmbracelet/bubbletea"

// Bubble Tea application state.
type model struct {
	game          *game // game state
	width, height int   // terminal dimensions

	CurrentWindow Window
	MenuWin       *MenuWindow
	BoardWin      *BoardWindow
}

func (model model) Init() tea.Cmd { return nil }

// Creates a new model with initial state.
func InitialModel() model {
	menuWinInst := NewMenuWindow()

	return model{
		game:          newGame(),
		CurrentWindow: menuWinInst,
		MenuWin:       menuWinInst,
		BoardWin:      NewBoardWindow(),
	}
}
