package main

// Bubble Tea application state.
type model struct {
	game          game // game state
	inGame        bool // whether currently playing (vs. in menu)
	menuChoice    int  // current menu selection
	width, height int  // terminal dimensions

	menuMinW, menuMinH   int // minimum terminal size for menu
	boardMinW, boardMinH int // minimum terminal size for game board
}

// Creates a new model with initial state.
func initialModel() model {
	board := make([][]cell, rows)
	for i := range board {
		board[i] = make([]cell, cols)
	}

	game := game{
		board:       board,
		cursorRow:   0,
		cursorCol:   0,
		numMines:    mines,
		state:       playing,
		minesPlaced: false,
	}

	return model{
		game:       game,
		inGame:     false,
		menuChoice: 0,

		menuMinW:  30,
		menuMinH:  10,
		boardMinW: 35,
		boardMinH: 15,
	}
}
