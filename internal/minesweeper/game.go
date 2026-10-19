package minesweeper

import (
	"hash/fnv"
	"math/rand"
	"time"
)

// CellState is a bitmask of
// what is true of a cell.
type CellState int

// States a cell can be in, combined in its CellState.
const (
	Revealed CellState = 1 << iota
	Flagged
	Mined
	Exploded
	Odd
)

// Cell is one square of the board.
type Cell struct {
	State CellState // bitmasked state of the cell
	Adj   int       // number of adjacent mines
}

// StateIs reports whether the cell is in any of the given states.
func (cell Cell) StateIs(state CellState) bool {
	return cell.State&state != 0
}

// GameState is a bitmask of
// where the game is at.
type GameState int

// States a game can be in,
// combined in its GameState.
const (
	Playing GameState = 1 << iota
	FirstMove
	Lost
	Won

	PlayerTurn = Playing | FirstMove
)

// Game state.
type Game struct {
	Timer
	settings Settings

	board [][]Cell
	state GameState // bitmasked state of the game
	rng   *rand.Rand

	usedFlags   int // flags used by the player
	numCells    int // total cells count
	numSafe     int // non-mine cells count
	numRevealed int // revealed cells count
}

// StateIs reports whether the game is in any of the given states.
func (game *Game) StateIs(state GameState) bool {
	return game.state&state != 0
}

// New initializes a new game with the given settings.
func New(settings Settings) *Game {
	settings.Rows = max(settings.Rows, 0)
	settings.Cols = max(settings.Cols, 0)
	settings.Mines = min(max(settings.Mines, 0), max(settings.MaxMines(), 0))

	board := make([][]Cell, settings.Rows)
	for i := range board {
		board[i] = make([]Cell, settings.Cols)

		for j := range board[i] {
			if (i+j)%2 == 1 {
				board[i][j].State |= Odd
			}
		}
	}

	numCells := settings.Rows * settings.Cols

	randSource := rand.NewSource(
		seedFrom(settings.Seed),
	)

	return &Game{
		Timer:    Timer{},
		settings: settings,
		board:    board,
		state:    FirstMove,
		rng:      rand.New(randSource),
		numCells: numCells,
		numSafe:  numCells - settings.Mines,
	}
}

// seedFrom hashes a seed string into a random source seed.
//
//	NOTE: a empty string returns a random seed.
func seedFrom(seed string) int64 {
	if seed == "" {
		return time.Now().UnixNano()
	}

	hash := fnv.New64a()
	hash.Write([]byte(seed))

	return int64(hash.Sum64())
}

// Rows is the board's height.
func (game *Game) Rows() int { return game.settings.Rows }

// Cols is the board's width.
func (game *Game) Cols() int { return game.settings.Cols }

// Cell returns the cell at the given coordinates.
func (game *Game) Cell(row, col int) Cell { return game.board[row][col] }

// checks if the coordinates are within the game board.
func (game *Game) inBounds(row, col int) bool {
	inRow := row >= 0 && row < game.settings.Rows
	inCol := col >= 0 && col < game.settings.Cols
	return inRow && inCol
}

// Calls fn for each valid neighbor of the cell.
func (game *Game) forEachNeighbour(row, col int, fn func(r, c int)) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			nr, nc := row+i, col+j
			if game.inBounds(nr, nc) {
				fn(nr, nc)
			}
		}
	}
}

// Randomly places mines on the board, avoiding the first revealed cell.
func (game *Game) placeMines(firstR, firstC int) {
	placed := 0
	for placed < game.settings.Mines {
		row := game.rng.Intn(game.settings.Rows)
		col := game.rng.Intn(game.settings.Cols)
		if game.board[row][col].StateIs(Mined) {
			continue
		}

		if row == firstR && col == firstC {
			continue
		}

		game.board[row][col].State |= Mined
		game.forEachNeighbour(row, col, func(nr, nc int) {
			if !game.board[nr][nc].StateIs(Mined) {
				game.board[nr][nc].Adj++
			}
		})

		placed++
	}
}

// Reveals all mine locations on the board
func (game *Game) revealAllMines() {
	for row := range game.board {
		for col := range game.board[row] {
			if game.board[row][col].StateIs(Mined) {
				game.board[row][col].State |= Revealed
			}
		}
	}
}

// Reveal uncovers the cell at the given coordinates, or its unflagged
// neighbours when it is revealed and flagged around as many times as its
// number. The mines are placed on the first move, so that cell is never
// mined. The game ends on a mine or once every safe cell is revealed.
func (game *Game) Reveal(row, col int) {
	if !game.StateIs(PlayerTurn) || !game.inBounds(row, col) {
		return
	}

	if game.StateIs(FirstMove) {
		game.placeMines(row, col)
		game.state = Playing
		game.StartTimer()
	}

	cell := &game.board[row][col]
	if cell.StateIs(Flagged) {
		return
	}

	if cell.StateIs(Revealed) {
		game.revealAround(row, col)
	} else {
		game.revealSingleCell(row, col)
	}

	game.checkWin()
}

// Reveals a single cell and reveals neighbors if it has no adjacent mines.
// Expects the cell to be in bounds and not revealed nor flagged.
func (game *Game) revealSingleCell(row int, col int) {
	pending := [][2]int{{row, col}}

	for len(pending) > 0 {
		last := len(pending) - 1
		row, col := pending[last][0], pending[last][1]
		pending = pending[:last]

		cell := &game.board[row][col]
		if cell.StateIs(Revealed) || cell.StateIs(Flagged) {
			continue
		}

		game.numRevealed++
		cell.State |= Revealed
		if cell.StateIs(Mined) {
			cell.State |= Exploded
			game.state = Lost
			game.StopTimer()
			game.revealAllMines()
			return
		}

		if cell.Adj != 0 {
			continue
		}

		// Queues empty neighboring cells.
		game.forEachNeighbour(row, col, func(nr, nc int) {
			cell := &game.board[nr][nc]
			if !cell.StateIs(Revealed) && !cell.StateIs(Flagged) {
				pending = append(pending, [2]int{nr, nc})
			}
		})
	}
}

// Reveals around a cell with adjacent flagged cells.
// Expects the cell to be in bounds, revealed and not flagged.
func (game *Game) revealAround(row, col int) {
	cell := &game.board[row][col]
	if cell.Adj <= 0 {
		return
	}

	numAdjFlags := 0
	game.forEachNeighbour(row, col, func(nr, nc int) {
		cell := &game.board[nr][nc]
		if cell.StateIs(Flagged) {
			numAdjFlags++
		}
	})

	if numAdjFlags < cell.Adj {
		return
	}

	game.forEachNeighbour(row, col, func(nr, nc int) {
		cell := &game.board[nr][nc]
		if !cell.StateIs(Flagged) && !cell.StateIs(Revealed) {
			game.revealSingleCell(nr, nc)
		}
	})
}

// Checks if all non-mine cells have been revealed.
func (game *Game) checkWin() {
	if !game.StateIs(Playing) {
		return
	}

	if game.numRevealed >= game.numSafe {
		game.state = Won
		game.StopTimer()
	}
}

// ToggleFlag flags or unflags the cell at the given coordinates.
func (game *Game) ToggleFlag(row, col int) {
	if !game.StateIs(PlayerTurn) || !game.inBounds(row, col) {
		return
	}

	cell := &game.board[row][col]
	if cell.StateIs(Revealed) {
		return
	}

	cell.State ^= Flagged
	if cell.StateIs(Flagged) {
		game.usedFlags++
		return
	}

	game.usedFlags--
}

// FlagsRemaining is the mine count less the flags placed.
func (game *Game) FlagsRemaining() int {
	return game.settings.Mines - game.usedFlags
}
