package board

import (
	"testing"

	"termsweeper/internal/config"
	"termsweeper/internal/minesweeper"
	"termsweeper/internal/tui/theme"
)

// newTestWindow builds a board window sized for the given game,
// with the default theme unless the caller adjusts it.
func newTestWindow(t *testing.T, settings minesweeper.Settings) *Window {
	t.Helper()

	window := New(theme.NewTheme(config.Default()))
	window.NewGame(settings)

	return window
}

func TestScrollOffset(t *testing.T) {
	tests := []struct {
		name      string
		offset    int
		cursor    int
		visible   int
		total     int
		scrolloff int
		want      int
	}{
		{"cursor already visible", 0, 3, 10, 30, 0, 0},
		{"cursor below the viewport", 0, 12, 10, 30, 0, 3},
		{"cursor above the viewport", 10, 4, 10, 30, 0, 4},
		{"scrolloff keeps rows below", 0, 8, 10, 30, 2, 1},
		{"scrolloff keeps rows above", 10, 11, 10, 30, 2, 9},
		{"clamped at the start", 0, 0, 10, 30, 3, 0},
		{"clamped at the end", 20, 29, 10, 30, 3, 20},
		{"scrolloff larger than the viewport centres", 0, 20, 5, 30, 100, 18},
		{"board shorter than the viewport", 0, 2, 10, 4, 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := scrollOffset(test.offset, test.cursor, test.visible, test.total, test.scrolloff)

			if got != test.want {
				t.Errorf("scrollOffset(%d, %d, %d, %d, %d) = %d, want %d",
					test.offset, test.cursor, test.visible, test.total, test.scrolloff, got, test.want)
			}
		})
	}
}

func TestViewportFitsTheChrome(t *testing.T) {
	window := newTestWindow(t, minesweeper.Settings{Rows: 40, Cols: 40, Mines: 10})

	rows, cols := window.viewport(80, 24)

	if want := 24 - window.chromeHeight(); rows != want {
		t.Errorf("viewport rows = %d, want %d", rows, want)
	}

	if want := (80 - window.chromeWidth()) / 2; cols != want {
		t.Errorf("viewport cols = %d, want %d", cols, want)
	}

	// A board smaller than the terminal is never padded out.
	window.NewGame(minesweeper.Settings{Rows: 5, Cols: 5, Mines: 1})

	if rows, cols := window.viewport(200, 200); rows != 5 || cols != 5 {
		t.Errorf("viewport = %dx%d, want 5x5 for a board that fits", rows, cols)
	}

	// A terminal too small to hold anything still leaves a cell to draw.
	if rows, cols := window.viewport(1, 1); rows != 1 || cols != 1 {
		t.Errorf("viewport = %dx%d, want 1x1 at minimum", rows, cols)
	}
}

func TestScrollToCursorFollowsTheCursor(t *testing.T) {
	window := newTestWindow(t, minesweeper.Settings{Rows: 40, Cols: 40, Mines: 10})
	width, height := 40, 20
	visibleRows, visibleCols := window.viewport(width, height)

	window.CursorRow = 39
	window.CursorCol = 39
	window.scrollToCursor(width, height)

	if got, want := window.offsetRow, 40-visibleRows; got != want {
		t.Errorf("offsetRow = %d, want %d", got, want)
	}

	if got, want := window.offsetCol, 40-visibleCols; got != want {
		t.Errorf("offsetCol = %d, want %d", got, want)
	}

	window.CursorRow = 0
	window.CursorCol = 0
	window.scrollToCursor(width, height)

	if window.offsetRow != 0 || window.offsetCol != 0 {
		t.Errorf("offset = (%d, %d), want (0, 0)", window.offsetRow, window.offsetCol)
	}
}

func TestViewRendersEveryVisibleRow(t *testing.T) {
	window := newTestWindow(t, minesweeper.Settings{Rows: 40, Cols: 40, Mines: 10})

	// Rendering a board larger than the terminal must not panic
	// or read outside the board.
	if window.View(40, 20) == "" {
		t.Fatal("View returned nothing")
	}

	window.CursorRow = 39
	window.CursorCol = 39

	if window.View(40, 20) == "" {
		t.Fatal("View returned nothing with the cursor at the far corner")
	}
}
