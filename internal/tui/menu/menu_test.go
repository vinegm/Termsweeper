package menu

import (
	"testing"

	"termsweeper/internal/config"
	"termsweeper/internal/minesweeper"
	"termsweeper/internal/tui/theme"
)

func newTestWindow(t *testing.T, settings minesweeper.Settings) *Window {
	t.Helper()

	return New(theme.NewTheme(config.Default()), settings)
}

// selectRow moves the menu cursor onto the row editing the given setting.
func selectRow(t *testing.T, window *Window, action menuOption) {
	t.Helper()

	for i, item := range window.items() {
		if item.action == action {
			window.Choice = i
			return
		}
	}

	t.Fatalf("no menu row edits %v", action)
}

func TestAdjustKeepsSettingsPlayable(t *testing.T) {
	window := newTestWindow(t, minesweeper.Settings{Rows: minesweeper.MinDim, Cols: minesweeper.MinDim, Mines: 1})

	window.adjust(boardRows, -1)
	window.adjust(boardCols, -1)

	if window.Settings.Rows != minesweeper.MinDim || window.Settings.Cols != minesweeper.MinDim {
		t.Errorf("board shrank to %dx%d, want the %d minimum held",
			window.Settings.Rows, window.Settings.Cols, minesweeper.MinDim)
	}

	window.adjust(mineCount, -1)

	if got, want := window.Settings.Mines, minesweeper.MinMines; got != want {
		t.Errorf("Mines = %d, want the %d minimum held", got, want)
	}

	window.Settings.Rows = minesweeper.MaxDim
	window.adjust(boardRows, 1)

	if got, want := window.Settings.Rows, minesweeper.MaxDim; got != want {
		t.Errorf("Rows = %d, want the %d maximum held", got, want)
	}

	// The mine count can never fill the board: the first move needs a safe cell.
	window.Settings = minesweeper.Settings{Rows: 5, Cols: 5, Mines: 24}
	window.adjust(mineCount, 1)

	if got, want := window.Settings.Mines, 24; got != want {
		t.Errorf("Mines = %d, want %d", got, want)
	}

	// Shrinking the board pulls an over-large mine count down with it.
	window.Settings = minesweeper.Settings{Rows: 6, Cols: 5, Mines: 29}
	window.adjust(boardRows, -1)

	if got, want := window.Settings.Mines, 24; got != want {
		t.Errorf("Mines = %d, want %d after the board shrank", got, want)
	}

	if err := window.Settings.Validate(); err != nil {
		t.Errorf("menu produced unplayable settings: %v", err)
	}
}

func TestHandleKeyNavigatesAndStarts(t *testing.T) {
	window := newTestWindow(t, minesweeper.Settings{Rows: 16, Cols: 16, Mines: 40})

	window.handleKey("up")

	if window.Choice != 0 {
		t.Errorf("Choice = %d, want it held at the first row", window.Choice)
	}

	for range len(window.items()) + 2 {
		window.handleKey("down")
	}

	if got, want := window.Choice, len(window.items())-1; got != want {
		t.Errorf("Choice = %d, want it held at the last row %d", got, want)
	}

	selectRow(t, window, boardCols)
	window.handleKey("right")

	if got, want := window.Settings.Cols, 17; got != want {
		t.Errorf("Cols = %d, want %d", got, want)
	}

	selectRow(t, window, startNewGame)
	cmd := window.handleKey("enter")

	if cmd == nil {
		t.Fatal("enter on New Game returned no command")
	}

	msg, ok := cmd().(StartGameMsg)
	if !ok {
		t.Fatalf("enter on New Game returned %T, want StartGameMsg", cmd())
	}

	if msg.Settings != window.Settings {
		t.Errorf("started with %+v, want %+v", msg.Settings, window.Settings)
	}
}

func TestContinueRowOnlyExistsDuringAGame(t *testing.T) {
	window := newTestWindow(t, minesweeper.Settings{Rows: 16, Cols: 16, Mines: 40})

	for _, item := range window.items() {
		if item.action == continueGame {
			t.Fatal("Continue is offered with no game in progress")
		}
	}

	window.OngoingGame = true
	selectRow(t, window, continueGame)

	cmd := window.handleKey("enter")
	if cmd == nil {
		t.Fatal("enter on Continue returned no command")
	}

	if _, ok := cmd().(ContinueGameMsg); !ok {
		t.Errorf("enter on Continue returned %T, want ContinueGameMsg", cmd())
	}
}
