package board

import (
	"fmt"
	"strings"

	"github.com/vinegm/termsweeper/internal/config"
)

// rowLabel renders the number shown on one side of a board row, full width so
// e.g. row 10 doesn't get clipped down to its last digit, or nothing when that
// side's row labels are turned off.
func (window *Window) rowLabel(row int, show bool) string {
	if !show {
		return ""
	}

	number := labelNumber(row, window.CursorRow, window.opts.RelativeLineNumbers)
	text := fmt.Sprintf("%*d ", window.rowNumberDigits(), number)

	return window.theme.LabelStyle(row, window.CursorRow).Render(text)
}

// columnLabels renders the number block shown above and below the grid, one
// digit per line with the least significant digit closest to the grid. Each
// column only has a single character-cell to spare, so multi-digit column
// numbers stack vertically there instead of overflowing sideways.
func (window *Window) columnLabels(visibleCols int) string {
	height := window.columnNumberDigits()
	leftPad := window.theme.TextStyle.Render(strings.Repeat(" ", window.labelPadding(window.opts.Numbers.Has(config.EdgeLeft))))
	rightPad := window.theme.TextStyle.Render(strings.Repeat(" ", window.labelPadding(window.opts.Numbers.Has(config.EdgeRight))))

	lines := make([]string, height)

	for row := range height {
		var sb strings.Builder
		sb.WriteString(leftPad)

		for col := window.offsetCol; col < window.offsetCol+visibleCols; col++ {
			number := labelNumber(col, window.CursorCol, window.opts.RelativeColumnNumbers)
			digit := fmt.Sprintf("%*d", height, number)[row]

			sb.WriteString(window.theme.LabelStyle(col, window.CursorCol).Render(string(digit) + " "))
		}

		sb.WriteString(rightPad)
		lines[row] = sb.String()
	}

	return strings.Join(lines, "\n")
}

// labelPadding is the width the column labels have to skip on one side to line up
// with the grid: the arrow gutter, plus the row label when that side shows one.
func (window *Window) labelPadding(showRowLabel bool) int {
	padding := 1 // arrow gutter

	if showRowLabel {
		padding += window.rowNumberDigits() + 1 // digits, plus the trailing space
	}

	return padding
}

// rowNumberDigits is how many decimal digits the widest row number or relative
// distance can take.
func (window *Window) rowNumberDigits() int {
	return digitWidth(window.game.Rows() - 1)
}

// columnNumberDigits is how many decimal digits the widest column number or
// relative distance can take.
func (window *Window) columnNumberDigits() int {
	return digitWidth(window.game.Cols() - 1)
}

// digitWidth returns how many decimal digits n takes, minimum 1.
func digitWidth(n int) int {
	width := 1
	for n >= 10 {
		n /= 10
		width++
	}

	return width
}

// labelNumber returns the number to display for a row/column label, either
// absolute or relative to the cursor's position depending on the relative flag.
func labelNumber(idx, cursorIdx int, relative bool) int {
	if relative && idx != cursorIdx {
		return abs(idx - cursorIdx)
	}

	return idx
}

func abs(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
