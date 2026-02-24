package board

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"termsweeper/internal/minesweeper"
	"termsweeper/internal/tui/theme"
)

// Space the board window needs on top of the board itself.
const (
	boardChromeWidth  = 6 // borders + padding + arrow gutters
	boardChromeHeight = 6 // title + hint + arrow gutters + borders
)

const boardTitle = "Termsweeper"

type (
	// Asks for control to exit the board window.
	ExitMsg struct{}

	// Message carrying a clock tick, tagged with the game generation that armed it.
	clockTickMsg struct{ tag int }
)

// Implements the game/board screen and holds its view state.
type Window struct {
	CursorRow int
	CursorCol int

	// Board cell shown at the top-left of the window.
	offsetRow int
	offsetCol int

	// vim-like movement count.
	countBuffer string

	theme    *theme.Theme
	game     *minesweeper.Game
	settings minesweeper.Settings

	clockTag int
	ticking  bool
}

// New initializes a new Window.
//
//	NOTE: The game is not started until newGame() is called.
func New(th *theme.Theme) *Window {
	return &Window{theme: th}
}

// NewGame Starts a fresh game from the given settings.
func (window *Window) NewGame(settings minesweeper.Settings) {
	window.CursorRow = 0
	window.CursorCol = 0
	window.offsetRow = 0
	window.offsetCol = 0
	window.countBuffer = ""
	window.settings = settings
	window.game = minesweeper.New(settings)

	window.clockTag++
	window.ticking = false
}

func (window *Window) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		cmd := window.handleKey(strings.ToLower(msg.String()))

		if cmd != nil {
			return cmd
		}

		return window.startClock()

	case clockTickMsg:
		if msg.tag != window.clockTag {
			return nil
		}

		if !window.game.StateIs(minesweeper.Playing) {
			window.ticking = false
			return nil
		}

		return clockTick(window.clockTag)
	}

	return nil
}

func (window *Window) View(width, height int) string {
	window.scrollToCursor(width, height)

	visibleRows, visibleCols := window.viewport(width, height)
	clipped := window.clippedEdges(visibleRows, visibleCols)

	var sb strings.Builder

	sb.WriteString(window.theme.TitleStyle.Render(boardTitle))
	sb.WriteString("\n")
	sb.WriteString(window.theme.TextStyle.Render(window.hint()))
	sb.WriteString("\n")
	sb.WriteString(window.gutter(clipped.top, theme.ArrowUp))
	sb.WriteString("\n")

	if window.theme.ColumnNumbers.Top() {
		sb.WriteString(window.columnLabels(visibleCols))
		sb.WriteString("\n")
	}

	// The side arrows sit halfway down the grid, next to the outermost columns.
	arrowRow := window.offsetRow + visibleRows/2

	for row := window.offsetRow; row < window.offsetRow+visibleRows; row++ {
		onArrowRow := row == arrowRow

		sb.WriteString(window.gutter(clipped.left && onArrowRow, theme.ArrowLeft))
		sb.WriteString(window.rowLabel(row, window.theme.LineNumbers.Left()))

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

		sb.WriteString(window.rowLabel(row, window.theme.LineNumbers.Right()))
		sb.WriteString(window.gutter(clipped.right && onArrowRow, theme.ArrowRight))
		sb.WriteString("\n")
	}

	if window.theme.ColumnNumbers.Bottom() {
		sb.WriteString(window.columnLabels(visibleCols))
		sb.WriteString("\n")
	}

	sb.WriteString(window.gutter(clipped.bottom, theme.ArrowDown))

	return window.theme.WindowStyle.Render(sb.String())
}

// Window edges the board runs past, i.e. the ones with more board behind them.
type clippedEdges struct {
	top, right, bottom, left bool
}

func (window *Window) clippedEdges(visibleRows, visibleCols int) clippedEdges {
	return clippedEdges{
		top:    window.offsetRow > 0,
		bottom: window.offsetRow+visibleRows < window.game.Rows(),
		left:   window.offsetCol > 0,
		right:  window.offsetCol+visibleCols < window.game.Cols(),
	}
}

// gutter renders one cell of gutter: the arrow when the board continues past that
// edge, a blank otherwise. Always one cell wide, so the grid never shifts as the
// arrows come and go.
func (window *Window) gutter(showArrow bool, arrow string) string {
	if !showArrow {
		arrow = " "
	}

	return window.theme.ArrowStyle.Render(arrow)
}

// rowLabel renders the number shown on one side of a board row, full width so
// e.g. row 10 doesn't get clipped down to its last digit, or nothing when that
// side's row labels are turned off.
func (window *Window) rowLabel(row int, show bool) string {
	if !show {
		return ""
	}

	number := labelNumber(row, window.CursorRow, window.theme.RelativeLineNumbers)
	text := fmt.Sprintf("%*d ", window.rowNumberDigits(), number)

	return window.theme.LabelStyle(row, window.CursorRow).Render(text)
}

// columnLabels renders the number block shown above and below the grid, one
// digit per line with the least significant digit closest to the grid. Each
// column only has a single character-cell to spare, so multi-digit column
// numbers stack vertically there instead of overflowing sideways.
func (window *Window) columnLabels(visibleCols int) string {
	height := window.columnNumberDigits()
	leftPad := window.theme.TextStyle.Render(strings.Repeat(" ", window.labelPadding(window.theme.LineNumbers.Left())))
	rightPad := window.theme.TextStyle.Render(strings.Repeat(" ", window.labelPadding(window.theme.LineNumbers.Right())))

	lines := make([]string, height)

	for row := range height {
		var sb strings.Builder
		sb.WriteString(leftPad)

		for col := window.offsetCol; col < window.offsetCol+visibleCols; col++ {
			number := labelNumber(col, window.CursorCol, window.theme.RelativeColumnNumbers)
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

// isHighlighted reports whether a cell falls in the cursor's highlighted row or
// column. The highlight is split at the cursor: HighlightRow lights up cells to
// its left and/or right, HighlightCol cells above and/or below, per their
// configured edge(s). The cursor's own cell is never highlighted here since it
// always renders with the selected-cell style instead.
func (window *Window) isHighlighted(row, col int) bool {
	if window.CursorRow == row {
		side := window.theme.HighlightRow
		if (col < window.CursorCol && side.Left()) || (col > window.CursorCol && side.Right()) {
			return true
		}
	}

	if window.CursorCol == col {
		side := window.theme.HighlightCol
		if (row < window.CursorRow && side.Top()) || (row > window.CursorRow && side.Bottom()) {
			return true
		}
	}

	return false
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

// chromeWidth is the horizontal space taken by everything but the grid itself.
func (window *Window) chromeWidth() int {
	width := boardChromeWidth

	if window.theme.LineNumbers.Left() {
		width += window.rowNumberDigits() + 1
	}

	if window.theme.LineNumbers.Right() {
		width += window.rowNumberDigits() + 1
	}

	return width
}

// chromeHeight is the vertical space taken by everything but the grid itself.
func (window *Window) chromeHeight() int {
	height := boardChromeHeight

	if window.theme.ColumnNumbers.Top() {
		height += window.columnNumberDigits()
	}

	if window.theme.ColumnNumbers.Bottom() {
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

// scrollOffset returns the new offset, which may be the same as the old one. off
// is the minimum number of cells kept visible between the cursor and the
// viewport edge, like nvim's 'scrolloff'/'sidescrolloff' — capped to half the
// viewport so it can't fight itself when the window is too small to honor it,
// and naturally ignored near either end of the board since offset stays clamped
// within [0, total-visible].
func scrollOffset(offset, cursor, visible, total, off int) int {
	off = clamp(off, 0, (visible-1)/2)

	if cursor-off < offset {
		offset = cursor - off
	}

	if cursor+off >= offset+visible {
		offset = cursor + off - visible + 1
	}

	return clamp(offset, 0, total-visible)
}

func clamp(value, low, high int) int {
	return min(max(value, low), high)
}

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
// length, defaulting to 1 when no count was entered. Movement wraps, so a full lap
// is a no-op and counts fold into the axis length — which also keeps a pasted-in
// count from spinning the cursor for millions of iterations.
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

func (window *Window) moveUp() {
	if window.CursorRow == 0 {
		window.CursorRow = window.game.Rows() - 1
		return
	}

	window.CursorRow--
}

func (window *Window) moveDown() {
	if window.CursorRow == window.game.Rows()-1 {
		window.CursorRow = 0
		return
	}

	window.CursorRow++
}

func (window *Window) moveLeft() {
	if window.CursorCol == 0 {
		window.CursorCol = window.game.Cols() - 1
		return
	}

	window.CursorCol--
}

func (window *Window) moveRight() {
	if window.CursorCol == window.game.Cols()-1 {
		window.CursorCol = 0
		return
	}

	window.CursorCol++
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
	case "up", "k":
		for range rowCount {
			window.moveUp()
		}

	case "down", "j":
		for range rowCount {
			window.moveDown()
		}

	case "left", "h":
		for range colCount {
			window.moveLeft()
		}

	case "right", "l":
		for range colCount {
			window.moveRight()
		}

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

func (window *Window) startClock() tea.Cmd {
	if window.ticking || !window.game.StateIs(minesweeper.Playing) {
		return nil
	}

	window.ticking = true

	return clockTick(window.clockTag)
}

func clockTick(tag int) tea.Cmd {
	notifyTick := func(time.Time) tea.Msg {
		return clockTickMsg{tag: tag}
	}

	return tea.Tick(time.Second, notifyTick)
}
