package board

import (
	"strconv"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/vinegm/termsweeper/internal/minesweeper"
)

// appendCount accumulates digits for a "count + direction" move, e.g. "5j" moves
// down 5 times. Reports whether the key was consumed as part of a count. A leading
// "0" doesn't start one, matching vim, so it stays free for other bindings.
func (window *Window) appendCount(key string) bool {
	if len(key) != 1 || key[0] < '0' || key[0] > '9' {
		return false
	}

	if key == "0" && window.countBuffer == "" {
		return false
	}

	window.countBuffer += key

	return true
}

// moveCount is how many times to repeat a movement along an axis of the given
// length, defaulting to 1 when no count was entered.
//
//	NOTE: Movement wraps.
func (window *Window) moveCount(axisLen int) int {
	count := 1

	if window.countBuffer != "" {
		parsed, err := strconv.Atoi(window.countBuffer)
		if err == nil && parsed > 0 {
			count = parsed
		}
	}

	if axisLen > 0 {
		count %= axisLen
	}

	return count
}

// move shifts the cursor by the given number of rows and columns
//
//	NOTE: Movement wraps.
func (window *Window) move(rowDelta, colDelta int) {
	rows := window.game.Rows()
	cols := window.game.Cols()

	if rows <= 0 || cols <= 0 {
		return
	}

	window.CursorRow = (window.CursorRow + rowDelta + rows) % rows
	window.CursorCol = (window.CursorCol + colDelta + cols) % cols
}

func (window *Window) handleKey(pressed lowered) tea.Cmd {
	if window.appendCount(string(pressed)) {
		return nil
	}

	game := window.game
	rowCount := window.moveCount(game.Rows())
	colCount := window.moveCount(game.Cols())
	window.countBuffer = ""

	switch {
	case key.Matches(pressed, keys.Up):
		window.move(-rowCount, 0)

	case key.Matches(pressed, keys.Down):
		window.move(rowCount, 0)

	case key.Matches(pressed, keys.Left):
		window.move(0, -colCount)

	case key.Matches(pressed, keys.Right):
		window.move(0, colCount)

	case key.Matches(pressed, keys.Reveal):
		game.Reveal(window.CursorRow, window.CursorCol)

	case key.Matches(pressed, keys.Flag):
		game.ToggleFlag(window.CursorRow, window.CursorCol)

	case key.Matches(pressed, keys.NewGame):
		window.NewGame(window.settings)

	case key.Matches(pressed, keys.Menu):
		return func() tea.Msg {
			return ExitMsg{}
		}
	}

	return nil
}

// ShortHelp is the keys for the key hints bar.
func (window *Window) ShortHelp() []key.Binding {
	if !window.game.StateIs(minesweeper.PlayerTurn) {
		return []key.Binding{keys.NewGame, keys.Menu}
	}

	return []key.Binding{keys.Move, keys.Reveal, keys.Flag, keys.NewGame, keys.Menu}
}
