package board

import (
	"strconv"

	tea "charm.land/bubbletea/v2"
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

func (window *Window) handleKey(key string) tea.Cmd {
	if window.appendCount(key) {
		return nil
	}

	game := window.game
	rowCount := window.moveCount(game.Rows())
	colCount := window.moveCount(game.Cols())
	window.countBuffer = ""

	switch key {
	case "up", "k", "w":
		window.move(-rowCount, 0)

	case "down", "j", "s":
		window.move(rowCount, 0)

	case "left", "h", "a":
		window.move(0, -colCount)

	case "right", "l", "d":
		window.move(0, colCount)

	case "space", "enter":
		game.Reveal(window.CursorRow, window.CursorCol)

	case "f":
		game.ToggleFlag(window.CursorRow, window.CursorCol)

	case "r":
		window.NewGame(window.settings)

	case "q":
		return func() tea.Msg {
			return ExitMsg{}
		}
	}

	return nil
}
