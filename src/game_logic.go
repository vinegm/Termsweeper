package src

import (
	"fmt"
	"math/rand"
)

/* Temporary constants for board dimensions and mine count.
 * In the complete implementation, these will be dinamic.
 */
const (
	rows  = 18
	cols  = 18
	mines = 40
)

// neighbor offsets for iterating over adjacent cells
var neighbors = [][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

type gameState int

const (
	playing gameState = iota
	lost
	won
)

// Board cell representation.
type cell struct {
	revealed bool // Has been revealed
	mined    bool // Has a mine
	flagged  bool // Is flagged
	adj      int  // number of adjacent mines
}

// Returns the char for a given cell based on its state.
func (cell *cell) char() string {
	if cell.revealed {
		if cell.mined {
			return AppConfig.MineChar
		}

		if cell.adj == 0 {
			return " "
		}

		return fmt.Sprintf("%d", cell.adj)
	}

	if cell.flagged {
		return AppConfig.FlagChar
	}

	return "·"
}

// Game state.
type game struct {
	board       [][]cell
	state       gameState // current game state
	minesPlaced bool      // have mines been placed
	numMines    int       // mines on the board

	usedFlags   int // flags used by the player
	numRevealed int // revealed cells count
}

func newGame() game {
	board := make([][]cell, rows)
	for i := range board {
		board[i] = make([]cell, cols)
	}

	return game{
		board:    board,
		numMines: mines,
		state:    playing,
	}
}

// inBounds checks if the given row and column coordinates are within the game board.
func (game *game) inBounds(row, col int) bool {
	inRow := row >= 0 && row < rows
	inCol := col >= 0 && col < cols
	return inRow && inCol
}

// Returns a slice of valid neighboring cell coordinates.
func (game *game) validNeighbours(row, col int) [][2]int {
	var valid [][2]int
	for _, offset := range neighbors {
		nr, nc := row+offset[0], col+offset[1]
		if game.inBounds(nr, nc) {
			valid = append(valid, [2]int{nr, nc})
		}
	}

	return valid
}

// Randomly places mines on the board, avoiding the first revealed cell.
func (game *game) placeMines(firstR, firstC int) {
	placed := 0
	for placed < mines {
		row := rand.Intn(rows)
		col := rand.Intn(cols)
		if game.board[row][col].mined {
			continue
		}

		if row == firstR && col == firstC {
			continue
		}

		game.board[row][col].mined = true
		for _, n := range game.validNeighbours(row, col) {
			nr, nc := n[0], n[1]
			if !game.board[nr][nc].mined {
				game.board[nr][nc].adj++
			}
		}

		placed++
	}
}

// Reveals all mine locations on the board
func (game *game) revealAllMines() {
	for row := range game.board {
		for col := range game.board[row] {
			if game.board[row][col].mined {
				game.board[row][col].revealed = true
			}
		}
	}
}

// Reveals the cell at the given coordinates.
func (game *game) reveal(row, col int) {
	if !game.inBounds(row, col) {
		return
	}

	cell := &game.board[row][col]
	if cell.revealed {
		game.revealAround(row, col)
		return
	}

	game.revealSingleCell(row, col)
}

// reveals the cell at the given coordinates. If the cell has no adjacent mines,
// it recursively reveals neighboring cells.
func (game *game) revealSingleCell(row int, col int) {
	if !game.inBounds(row, col) {
		return
	}

	cell := &game.board[row][col]
	if cell.flagged || cell.revealed {
		return
	}

	game.numRevealed++
	cell.revealed = true
	if cell.mined {
		game.state = lost
		game.revealAllMines()
		return
	}

	if cell.adj != 0 {
		return
	}

	for _, n := range game.validNeighbours(row, col) {
		nr, nc := n[0], n[1]
		if !game.board[nr][nc].revealed {
			game.revealSingleCell(nr, nc)
		}
	}
}

// Reveals around a cell with adjacent flagged cells.
func (game *game) revealAround(row, col int) {
	if !game.inBounds(row, col) {
		return
	}

	cell := &game.board[row][col]
	if !cell.revealed || cell.adj <= 0 {
		return
	}

	numAdjFlags := 0
	for _, n := range game.validNeighbours(row, col) {
		nr, nc := n[0], n[1]
		cell := &game.board[nr][nc]
		if cell.flagged {
			numAdjFlags++
		}
	}

	if numAdjFlags < cell.adj {
		return
	}

	for _, n := range game.validNeighbours(row, col) {
		nr, nc := n[0], n[1]
		cell := &game.board[nr][nc]
		if !cell.flagged || !cell.revealed {
			game.revealSingleCell(nr, nc)
		}
	}
}

// Checks if all non-mine cells have been revealed.
func (game *game) checkWin() {
	if game.numRevealed >= cols*rows-game.numMines {
		game.state = won
	}
}

// Toggles the flag of the cell at the given coords.
func (game *game) toggleFlag(row, col int) {
	if !game.inBounds(row, col) {
		return
	}

	cell := &game.board[row][col]
	if cell.revealed {
		return
	}

	if cell.flagged {
		cell.flagged = false
		game.usedFlags--
		return
	}

	cell.flagged = true
	game.usedFlags++
}

// Number of flags remaining.
func (game *game) flagsRemaining() int {
	return game.numMines - game.usedFlags
}
