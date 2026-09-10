package board

import (
	"termsweeper/internal/tui/theme"
)

// chromeWidth is the horizontal space taken by everything but the grid itself.
func (window *Window) chromeWidth() int {
	width := boardChromeWidth

	if window.theme.Numbers.Has(theme.EdgeLeft) {
		width += window.rowNumberDigits() + 1
	}

	if window.theme.Numbers.Has(theme.EdgeRight) {
		width += window.rowNumberDigits() + 1
	}

	return width
}

// chromeHeight is the vertical space taken by everything but the grid itself.
func (window *Window) chromeHeight() int {
	height := boardChromeHeight

	if window.theme.Numbers.Has(theme.EdgeTop) {
		height += window.columnNumberDigits()
	}

	if window.theme.Numbers.Has(theme.EdgeBottom) {
		height += window.columnNumberDigits()
	}

	return height
}

// viewport returns the number of rows and columns visible.
func (window *Window) viewport(width, height int) (rows, cols int) {
	rows = clamp(height-window.chromeHeight(), 1, window.game.Rows())
	cols = clamp((width-window.chromeWidth())/2, 1, window.game.Cols())

	return rows, cols
}

// scrollToCursor moves the viewport to keep the cursor visible.
func (window *Window) scrollToCursor(width, height int) {
	if window.game == nil {
		return
	}

	visibleRows, visibleCols := window.viewport(width, height)

	window.offsetRow = scrollOffset(window.offsetRow, window.CursorRow, visibleRows, window.game.Rows(), window.theme.Scrolloff)
	window.offsetCol = scrollOffset(window.offsetCol, window.CursorCol, visibleCols, window.game.Cols(), window.theme.SideScrolloff)
}

// scrollOffset returns the new offset, which may be the same as the old one.
func scrollOffset(offset, cursor, visible, total, off int) int {
	off = clamp(off, 0, (visible-1)/2)

	if cursor-off < offset {
		offset = cursor - off
	}

	if cursor+off >= offset+visible {
		offset = cursor + off - visible + 1
	}

	return clamp(offset, 0, max(total-visible, 0))
}

// clamp returns value clamped to the range [low, high].
func clamp(value, low, high int) int {
	return min(max(value, low), high)
}
