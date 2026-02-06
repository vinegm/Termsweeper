package src

import (
	"fmt"
	"math/rand"
)

// neighbor offsets
var neighbors = [][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

type cellState int

const (
	revealed cellState = 1 << iota
	flagged
	mined
	odd
)

// Board cell representation.
type cell struct {
	state cellState // bitmasked state of the cell
	adj   int       // number of adjacent mines
}

func (c *cell) is(state cellState) bool {
	return c.state&state != 0
}

// Returns the char for a given cell based on its state.
func (cell *cell) char() string {
	if cell.is(revealed) {
		if cell.is(mined) {
			return AppConfig.MineChar
		}

		if cell.adj == 0 {
			return " "
		}

		return fmt.Sprintf("%d", cell.adj)
	}

	if cell.is(flagged) {
		return AppConfig.FlagChar
	}

	return "·"
}

type gameState int

const (
	playing gameState = 1 << iota
	firstMove
	lost
	won

	playerTurn = playing | firstMove
)

// Game state.
type game struct {
	board [][]cell
	rows  int
	cols  int
	state gameState // bitmasked state of the game

	numMines    int // mines on the board
	usedFlags   int // flags used by the player
	numRevealed int // revealed cells count
}

func (game *game) is(state gameState) bool {
	return game.state&state != 0
}

func newGame() *game {
	numRows := AppConfig.BoardRows
	numCols := AppConfig.BoardCols
	numMines := AppConfig.MineCount

	board := make([][]cell, numRows)
	for i := range board {
		board[i] = make([]cell, numCols)

		for j := range board[i] {
			if (i+j)%2 == 1 {
				board[i][j].state |= odd
			}
		}
	}

	return &game{
		board:    board,
		rows:     numRows,
		cols:     numCols,
		numMines: numMines,
		state:    firstMove,
	}
}

// checks if the coordinates are within the game board.
func (game *game) inBounds(row, col int) bool {
	inRow := row >= 0 && row < game.rows
	inCol := col >= 0 && col < game.cols
	return inRow && inCol
}

// Calls fn for each valid neighbor of the cell.
func (game *game) forEachNeighbour(row, col int, fn func(r, c int)) {
	for _, offset := range neighbors {
		nr, nc := row+offset[0], col+offset[1]
		if game.inBounds(nr, nc) {
			fn(nr, nc)
		}
	}
}

// Randomly places mines on the board, avoiding the first revealed cell.
func (game *game) placeMines(firstR, firstC int) {
	placed := 0
	for placed < game.numMines {
		row := rand.Intn(game.rows)
		col := rand.Intn(game.cols)
		if game.board[row][col].is(mined) {
			continue
		}

		if row == firstR && col == firstC {
			continue
		}

		game.board[row][col].state |= mined
		game.forEachNeighbour(row, col, func(nr, nc int) {
			if !game.board[nr][nc].is(mined) {
				game.board[nr][nc].adj++
			}
		})

		placed++
	}
}

// Reveals all mine locations on the board
func (game *game) revealAllMines() {
	for row := range game.board {
		for col := range game.board[row] {
			if game.board[row][col].is(mined) {
				game.board[row][col].state |= revealed
			}
		}
	}
}

// Reveals the cell at the given coordinates.
// handles the first move logic and win/loss conditions.
// and calls the appropriate reveal functions based on the cell state.
func (game *game) reveal(row, col int) {
	if !game.is(playerTurn) || !game.inBounds(row, col) {
		return
	}

	if game.state == firstMove {
		game.placeMines(row, col)
		game.state = playing
	}

	cell := &game.board[row][col]
	if cell.is(flagged) {
		return
	}

	if cell.is(revealed) {
		game.revealAround(row, col)
	} else {
		game.revealSingleCell(row, col)
	}

	game.checkWin()
}

// Reveals a single cell and reveals neighbors if it has no adjacent mines.
// Expects the cell to be in bounds and not revealed nor flagged.
func (game *game) revealSingleCell(row int, col int) {
	cell := &game.board[row][col]

	game.numRevealed++
	cell.state |= revealed
	if cell.is(mined) {
		game.state = lost
		game.revealAllMines()
		return
	}

	if cell.adj != 0 {
		return
	}

	// Reveals empty neighboring cells.
	game.forEachNeighbour(row, col, func(nr, nc int) {
		cell := &game.board[nr][nc]
		if !cell.is(revealed) && !cell.is(flagged) {
			game.revealSingleCell(nr, nc)
		}
	})
}

// Reveals around a cell with adjacent flagged cells.
// Expects the cell to be in bounds, revealed and not flagged.
func (game *game) revealAround(row, col int) {
	cell := &game.board[row][col]
	if cell.adj <= 0 {
		return
	}

	numAdjFlags := 0
	game.forEachNeighbour(row, col, func(nr, nc int) {
		cell := &game.board[nr][nc]
		if cell.is(flagged) {
			numAdjFlags++
		}
	})

	if numAdjFlags < cell.adj {
		return
	}

	game.forEachNeighbour(row, col, func(nr, nc int) {
		cell := &game.board[nr][nc]
		if !game.inBounds(nr, nc) {
			return
		}

		if !cell.is(flagged) && !cell.is(revealed) {
			game.revealSingleCell(nr, nc)
		}
	})
}

// Checks if all non-mine cells have been revealed.
func (game *game) checkWin() {
	if game.numRevealed >= game.cols*game.rows-game.numMines {
		game.state = won
	}
}

// Toggles the flag of the cell at the given coords.
func (game *game) toggleFlag(row, col int) {
	if !game.is(playerTurn) || !game.inBounds(row, col) {
		return
	}

	cell := &game.board[row][col]
	if cell.is(revealed) {
		return
	}

	cell.state ^= flagged
	if cell.is(flagged) {
		game.usedFlags++
		return
	}

	game.usedFlags--
}

// Number of flags remaining.
func (game *game) flagsRemaining() int {
	return game.numMines - game.usedFlags
}
