package src

import tea "github.com/charmbracelet/bubbletea"

// Bubble Tea application state.
type model struct {
	game          game // game state (domain)
	inGame        bool // whether currently playing (vs. in menu)
	width, height int  // terminal dimensions

	CurrentWindow Window
	MenuWin       *MenuWindow
	BoardWin      *BoardWindow
}

func (model model) Init() tea.Cmd { return nil }

// Creates a new model with initial state.
func InitialModel() model {
	board := make([][]cell, rows)
	for i := range board {
		board[i] = make([]cell, cols)
	}

	game := game{
		board:       board,
		numMines:    mines,
		state:       playing,
		minesPlaced: false,
	}

	menuWinInst := &MenuWindow{Choice: 0, minWidth: 30, minHeight: 10}
	boardWinInst := &BoardWindow{minWidth: 35, minHeight: 15}

	return model{
		game:          game,
		inGame:        false,
		CurrentWindow: menuWinInst,
		MenuWin:       menuWinInst,
		BoardWin:      boardWinInst,
	}
}
