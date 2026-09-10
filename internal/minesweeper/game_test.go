package minesweeper

import "testing"

// newTestGame builds a game with a fixed mine layout, skipping the random
// placement the first move would do. '*' is a mine, any other rune a safe cell.
func newTestGame(t *testing.T, layout ...string) *Game {
	t.Helper()

	if len(layout) == 0 {
		t.Fatal("layout must have at least one row")
	}

	game := New(Settings{Rows: len(layout), Cols: len(layout[0])})
	game.state = Playing

	mines := 0

	for row, line := range layout {
		if len(line) != len(layout[0]) {
			t.Fatalf("row %d has %d cells, want %d", row, len(line), len(layout[0]))
		}

		for col, char := range line {
			if char == '*' {
				game.board[row][col].State |= Mined
				mines++
			}
		}
	}

	game.settings.Mines = mines
	game.numSafe = game.numCells - mines

	for row := range game.board {
		for col := range game.board[row] {
			if game.board[row][col].StateIs(Mined) {
				continue
			}

			game.forEachNeighbour(row, col, func(nr, nc int) {
				if game.board[nr][nc].StateIs(Mined) {
					game.board[row][col].Adj++
				}
			})
		}
	}

	return game
}

func countRevealed(game *Game) int {
	count := 0

	for row := range game.board {
		for col := range game.board[row] {
			if game.board[row][col].StateIs(Revealed) {
				count++
			}
		}
	}

	return count
}

func TestRevealFloodsEmptyRegion(t *testing.T) {
	game := newTestGame(t,
		".....",
		".....",
		".....",
		"....*",
	)

	game.Reveal(0, 0)

	// Everything but the mine opens: the only cells with adjacent
	// mines sit next to it and stop the flood at themselves.
	if got, want := countRevealed(game), 19; got != want {
		t.Fatalf("revealed %d cells, want %d", got, want)
	}

	if game.Cell(3, 4).StateIs(Revealed) {
		t.Error("mine was revealed by the flood fill")
	}

	if !game.StateIs(Won) {
		t.Errorf("state = %v, want Won after every safe cell opened", game.state)
	}
}

func TestRevealStopsAtNumberedCell(t *testing.T) {
	game := newTestGame(t,
		"*....",
		".....",
		".....",
	)

	game.Reveal(1, 1)

	if got := countRevealed(game); got != 1 {
		t.Fatalf("revealed %d cells, want 1 (cell has an adjacent mine)", got)
	}

	if got := game.Cell(1, 1).Adj; got != 1 {
		t.Errorf("Adj = %d, want 1", got)
	}
}

func TestRevealMineLoses(t *testing.T) {
	game := newTestGame(t,
		"*....",
		"....*",
		".....",
	)

	game.Reveal(0, 0)

	if !game.StateIs(Lost) {
		t.Fatalf("state = %v, want Lost", game.state)
	}

	if !game.Cell(0, 0).StateIs(Exploded) {
		t.Error("revealed mine is not marked Exploded")
	}

	if !game.Cell(1, 4).StateIs(Revealed) {
		t.Error("the other mine was not revealed on loss")
	}

	// The board is frozen after a loss.
	game.Reveal(2, 2)

	if game.Cell(2, 2).StateIs(Revealed) {
		t.Error("cell was revealed after the game was lost")
	}
}

func TestFlagBlocksRevealAndCountsDown(t *testing.T) {
	game := newTestGame(t,
		".....",
		"..*..",
		".....",
	)

	game.ToggleFlag(0, 0)

	if got, want := game.FlagsRemaining(), 0; got != want {
		t.Fatalf("FlagsRemaining = %d, want %d", got, want)
	}

	game.Reveal(0, 0)

	if game.Cell(0, 0).StateIs(Revealed) {
		t.Error("flagged cell was revealed")
	}

	game.ToggleFlag(0, 0)

	if got, want := game.FlagsRemaining(), 1; got != want {
		t.Fatalf("FlagsRemaining = %d, want %d after unflagging", got, want)
	}

	// Revealed cells can't be flagged.
	game.Reveal(0, 0)
	game.ToggleFlag(0, 0)

	if game.Cell(0, 0).StateIs(Flagged) {
		t.Error("revealed cell was flagged")
	}
}

func TestRevealAroundNeedsEnoughFlags(t *testing.T) {
	game := newTestGame(t,
		"*....",
		".....",
		".....",
	)

	game.Reveal(1, 1) // numbered cell, one adjacent mine
	game.Reveal(1, 1) // chord attempt with no flags placed

	if got := countRevealed(game); got != 1 {
		t.Fatalf("revealed %d cells, want 1: chording under-flagged should do nothing", got)
	}

	game.ToggleFlag(0, 0)
	game.Reveal(1, 1)

	if game.Cell(0, 0).StateIs(Revealed) {
		t.Error("chord revealed the flagged mine")
	}

	if !game.Cell(2, 2).StateIs(Revealed) {
		t.Error("chord did not open the surrounding cells")
	}
}

func TestFirstRevealIsNeverMined(t *testing.T) {
	// A board with a single safe cell: the first move has to land on it.
	settings := Settings{Rows: 5, Cols: 5, Mines: 24}

	for range 20 {
		game := New(settings)
		game.Reveal(2, 2)

		if game.StateIs(Lost) {
			t.Fatal("first revealed cell was mined")
		}

		if !game.StateIs(Won) {
			t.Fatalf("state = %v, want Won: the only safe cell was revealed", game.state)
		}
	}
}

func TestNewClampsUnplayableSettings(t *testing.T) {
	// More mines than cells would otherwise spin forever in placeMines.
	game := New(Settings{Rows: 4, Cols: 4, Mines: 100})
	game.Reveal(0, 0)

	if got, want := game.settings.Mines, 15; got != want {
		t.Errorf("Mines = %d, want %d", got, want)
	}

	if game.StateIs(FirstMove) {
		t.Error("reveal did not start the game")
	}

	empty := New(Settings{Rows: -3, Cols: -3, Mines: -3})

	if got := empty.Rows() + empty.Cols() + empty.settings.Mines; got != 0 {
		t.Errorf("negative settings were not clamped to zero, got %d", got)
	}
}
