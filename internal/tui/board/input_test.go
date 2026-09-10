package board

import (
	"testing"

	"termsweeper/internal/minesweeper"
)

func TestAppendCount(test *testing.T) {
	window := newTestWindow(test, minesweeper.Settings{Rows: 10, Cols: 10, Mines: 5})

	if window.appendCount("0") {
		test.Error("a leading 0 started a count, want it left free for other bindings")
	}

	for _, key := range []string{"1", "2"} {
		if !window.appendCount(key) {
			test.Fatalf("appendCount(%q) = false, want true", key)
		}
	}

	if !window.appendCount("0") {
		test.Error("appendCount(\"0\") = false, want true once a count has started")
	}

	if got, want := window.countBuffer, "120"; got != want {
		test.Errorf("countBuffer = %q, want %q", got, want)
	}

	if window.appendCount("j") {
		test.Error("appendCount consumed a non-digit key")
	}
}

func TestMoveCountWrapsWithinTheAxis(test *testing.T) {
	window := newTestWindow(test, minesweeper.Settings{Rows: 10, Cols: 10, Mines: 5})

	if got, want := window.moveCount(10), 1; got != want {
		test.Errorf("moveCount with no count = %d, want %d", got, want)
	}

	window.countBuffer = "3"

	if got, want := window.moveCount(10), 3; got != want {
		test.Errorf("moveCount = %d, want %d", got, want)
	}

	window.countBuffer = "23"

	if got, want := window.moveCount(10), 3; got != want {
		test.Errorf("moveCount = %d, want %d: counts wrap around the axis", got, want)
	}
}

func TestHandleKeyMovesTheCursor(test *testing.T) {
	window := newTestWindow(test, minesweeper.Settings{Rows: 10, Cols: 10, Mines: 5})

	for _, key := range []string{"5", "j"} {
		window.handleKey(key)
	}

	if got, want := window.CursorRow, 5; got != want {
		test.Fatalf("CursorRow = %d, want %d", got, want)
	}

	if window.countBuffer != "" {
		test.Errorf("countBuffer = %q, want it cleared after a move", window.countBuffer)
	}

	// Movement wraps in both directions.
	window.handleKey("k")
	window.handleKey("k")
	window.handleKey("k")
	window.handleKey("k")
	window.handleKey("k")
	window.handleKey("k")

	if got, want := window.CursorRow, 9; got != want {
		test.Errorf("CursorRow = %d, want %d after wrapping past the top", got, want)
	}

	window.handleKey("h")

	if got, want := window.CursorCol, 9; got != want {
		test.Errorf("CursorCol = %d, want %d after wrapping past the left edge", got, want)
	}
}

func TestHandleKeyPlaysTheGame(test *testing.T) {
	window := newTestWindow(test, minesweeper.Settings{Rows: 10, Cols: 10, Mines: 5})

	window.handleKey("f")

	if !window.game.Cell(0, 0).StateIs(minesweeper.Flagged) {
		test.Error("f did not flag the cell under the cursor")
	}

	window.handleKey("f")
	window.handleKey("space")

	if !window.game.Cell(0, 0).StateIs(minesweeper.Revealed) {
		test.Error("space did not reveal the cell under the cursor")
	}

	window.CursorRow = 4
	window.handleKey("r")

	if window.CursorRow != 0 || window.game.Cell(0, 0).StateIs(minesweeper.Revealed) {
		test.Error("r did not reset the board and the cursor")
	}

	if cmd := window.handleKey("q"); cmd == nil {
		test.Fatal("q returned no command, want an exit message")
	} else if _, ok := cmd().(ExitMsg); !ok {
		test.Error("q did not return an ExitMsg")
	}
}
