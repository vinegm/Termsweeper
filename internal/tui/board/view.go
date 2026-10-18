package board

import (
	"fmt"
	"strings"

	"github.com/vinegm/termsweeper/internal/config"
	"github.com/vinegm/termsweeper/internal/minesweeper"
	"github.com/vinegm/termsweeper/internal/tui/theme"
)

func (window *Window) View(width, height int) string {
	window.scrollToCursor(width, height)

	visibleRows, visibleCols := window.viewport(width, height)
	clipped := window.clippedEdges(visibleRows, visibleCols)

	var sb strings.Builder

	sb.WriteString(window.theme.TitleStyle.Render(boardTitle))
	sb.WriteString("\n")
	sb.WriteString(window.theme.TextStyle.Render(window.hint()))
	sb.WriteString("\n")
	sb.WriteString(window.gutter(clipped.Has(config.EdgeTop), theme.ArrowUp))
	sb.WriteString("\n")

	if window.opts.Numbers.Has(config.EdgeTop) {
		sb.WriteString(window.columnLabels(visibleCols))
		sb.WriteString("\n")
	}

	// The side arrows sit halfway down the grid, next to the outermost columns.
	arrowRow := window.offsetRow + visibleRows/2

	for row := window.offsetRow; row < window.offsetRow+visibleRows; row++ {
		onArrowRow := row == arrowRow

		sb.WriteString(window.gutter(clipped.Has(config.EdgeLeft) && onArrowRow, theme.ArrowLeft))
		sb.WriteString(window.rowLabel(row, window.opts.Numbers.Has(config.EdgeLeft)))

		for col := window.offsetCol; col < window.offsetCol+visibleCols; col++ {
			cell := window.game.Cell(row, col)
			char := window.theme.CellChar(cell)

			isCursorRow := window.CursorRow == row
			isCursorCol := window.CursorCol == col
			isSelectedCell := isCursorRow && isCursorCol
			isHighlighted := window.isHighlighted(row, col)
			style := window.theme.CellStyle(cell, isSelectedCell, isHighlighted)

			sb.WriteString(style.Render(char + " "))
		}

		sb.WriteString(window.rowLabel(row, window.opts.Numbers.Has(config.EdgeRight)))
		sb.WriteString(window.gutter(clipped.Has(config.EdgeRight) && onArrowRow, theme.ArrowRight))
		sb.WriteString("\n")
	}

	if window.opts.Numbers.Has(config.EdgeBottom) {
		sb.WriteString(window.columnLabels(visibleCols))
		sb.WriteString("\n")
	}

	sb.WriteString(window.gutter(clipped.Has(config.EdgeBottom), theme.ArrowDown))

	return window.theme.WindowStyle.Render(sb.String())
}

// clippedEdges reports which edges of the board the window clips.
func (window *Window) clippedEdges(visibleRows, visibleCols int) config.EdgeSetting {
	var clipped config.EdgeSetting

	if window.offsetRow > 0 {
		clipped |= config.EdgeTop
	}

	if window.offsetRow+visibleRows < window.game.Rows() {
		clipped |= config.EdgeBottom
	}

	if window.offsetCol > 0 {
		clipped |= config.EdgeLeft
	}

	if window.offsetCol+visibleCols < window.game.Cols() {
		clipped |= config.EdgeRight
	}

	return clipped
}

// gutter renders one cell of gutter: the arrow when the
// board continues past that edge, a blank otherwise.
//
//	WARN: Must be one cell wide, so the grid never shifts.
func (window *Window) gutter(showArrow bool, arrow string) string {
	if !showArrow {
		arrow = " "
	}

	return window.theme.ArrowStyle.Render(arrow)
}

// isHighlighted reports whether a cell falls in
// the cursor's highlighted row or column.
func (window *Window) isHighlighted(row, col int) bool {
	side := window.opts.Highlight
	onCursorRow := row == window.CursorRow
	onCursorCol := col == window.CursorCol

	switch {
	case onCursorRow && col < window.CursorCol:
		return side.Has(config.EdgeLeft)

	case onCursorRow && col > window.CursorCol:
		return side.Has(config.EdgeRight)

	case onCursorCol && row < window.CursorRow:
		return side.Has(config.EdgeTop)

	case onCursorCol && row > window.CursorRow:
		return side.Has(config.EdgeBottom)

	default:
		return false
	}
}

// hint returns the string to display
// in the window's hint line.
func (window *Window) hint() string {
	switch {
	case window.game.StateIs(minesweeper.Won):
		return fmt.Sprintf("You won! Time: %s", window.game.FormattedTime())

	case window.game.StateIs(minesweeper.Lost):
		return fmt.Sprintf("You lost! Time: %s", window.game.FormattedTime())

	default:
		return fmt.Sprintf("Flags: %d  Time: %s", window.game.FlagsRemaining(), window.game.FormattedTime())
	}
}
