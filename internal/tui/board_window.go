package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"

	"termsweeper/internal/minesweeper"
)

// Space the board window needs on top of the board itself.
const (
	boardChromeWidth  = 6 // borders + padding + arrow gutters
	boardChromeHeight = 6 // title + hint + arrow gutters + borders
)

const boardTitle = "Termsweeper"

// Board window messages.
type (
	// Asks for control to exit the board window.
	exitBoardMsg struct{}

	// Message carrying a clock tick, tagged with the game generation that armed it.
	clockTickMsg struct{ tag int }
)

// Implements the game/board screen and holds its view state.
type BoardWindow struct {
	CursorRow int
	CursorCol int

	// Board cell shown at the top-left of the window.
	offsetRow int
	offsetCol int

	theme    *Theme
	game     *minesweeper.Game
	settings minesweeper.Settings

	width, height int

	clockTag int
	ticking  bool
}

// NewBoardWindow initializes a new BoardWindow.
//
//	NOTE: The game is not started until newGame() is called.
func NewBoardWindow(theme *Theme) *BoardWindow {
	return &BoardWindow{theme: theme}
}

// NewGame Starts a fresh game from the given settings.
func (window *BoardWindow) NewGame(settings minesweeper.Settings) {
	window.CursorRow = 0
	window.CursorCol = 0
	window.offsetRow = 0
	window.offsetCol = 0
	window.settings = settings
	window.game = minesweeper.New(settings)

	window.clockTag++
	window.ticking = false
}

func (window *BoardWindow) Update(msg tea.Msg) (Window, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		window.width = msg.Width
		window.height = msg.Height
		window.scrollToCursor()

	case tea.KeyPressMsg:
		cmd := window.handleKey(strings.ToLower(msg.String()))
		window.scrollToCursor()

		if cmd != nil {
			return window, cmd
		}

		return window, window.startClock()

	case clockTickMsg:
		if msg.tag != window.clockTag {
			return window, nil
		}

		if !window.game.StateIs(minesweeper.Playing) {
			window.ticking = false
			return window, nil
		}

		return window, clockTick(window.clockTag)
	}

	return window, nil
}

func (window *BoardWindow) View() string {
	visibleRows, visibleCols := window.viewport()
	clipped := window.clippedEdges(visibleRows, visibleCols)

	var sb strings.Builder

	sb.WriteString(window.theme.titleStyle.Render(boardTitle))
	sb.WriteString("\n")
	sb.WriteString(window.theme.textStyle.Render(window.hint()))
	sb.WriteString("\n")
	sb.WriteString(window.gutter(clipped.top, arrowUp))
	sb.WriteString("\n")

	// The side arrows sit halfway down the grid, next to the outermost columns.
	arrowRow := window.offsetRow + visibleRows/2

	for row := window.offsetRow; row < window.offsetRow+visibleRows; row++ {
		onArrowRow := row == arrowRow

		sb.WriteString(window.gutter(clipped.left && onArrowRow, arrowLeft))

		for col := window.offsetCol; col < window.offsetCol+visibleCols; col++ {
			cell := window.game.Cell(row, col)
			char := window.theme.cellChar(cell)

			isCursorRow := window.CursorRow == row
			isCursorCol := window.CursorCol == col
			isSelectedCell := isCursorRow && isCursorCol
			style := window.theme.cellStyle(cell, isSelectedCell)

			sb.WriteString(style.Render(char + " "))
		}

		sb.WriteString(window.gutter(clipped.right && onArrowRow, arrowRight))
		sb.WriteString("\n")
	}

	sb.WriteString(window.gutter(clipped.bottom, arrowDown))

	return window.theme.windowStyle.Render(sb.String())
}

// Window edges the board runs past, i.e. the ones with more board behind them.
type clippedEdges struct {
	top, right, bottom, left bool
}

func (window *BoardWindow) clippedEdges(visibleRows, visibleCols int) clippedEdges {
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
func (window *BoardWindow) gutter(showArrow bool, arrow string) string {
	if !showArrow {
		arrow = " "
	}

	return window.theme.arrowStyle.Render(arrow)
}

func (window *BoardWindow) MinSize() (int, int) {
	contentWidth := max(2, lg.Width(boardTitle), lg.Width(window.hint()))

	return contentWidth + boardChromeWidth, 1 + boardChromeHeight
}

// viewport returns the number of rows and columns visible.
func (window *BoardWindow) viewport() (rows, cols int) {
	rows = clamp(window.height-boardChromeHeight, 1, window.game.Rows())
	cols = clamp((window.width-boardChromeWidth)/2, 1, window.game.Cols())

	return rows, cols
}

// scrollToCursor moves the viewport to keep the cursor visible.
func (window *BoardWindow) scrollToCursor() {
	if window.game == nil {
		return
	}

	visibleRows, visibleCols := window.viewport()

	window.offsetRow = scrollOffset(window.offsetRow, window.CursorRow, visibleRows, window.game.Rows())
	window.offsetCol = scrollOffset(window.offsetCol, window.CursorCol, visibleCols, window.game.Cols())
}

// scrollOffset returns the new offset, which may be the same as the old one.
func scrollOffset(offset, cursor, visible, total int) int {
	if cursor < offset {
		offset = cursor
	}

	if cursor >= offset+visible {
		offset = cursor - visible + 1
	}

	return clamp(offset, 0, total-visible)
}

func clamp(value, low, high int) int {
	return min(max(value, low), high)
}

func (window *BoardWindow) hint() string {
	switch {
	case window.game.StateIs(minesweeper.Won):
		return fmt.Sprintf("You won! Time: %s", window.game.FormattedTime())

	case window.game.StateIs(minesweeper.Lost):
		return fmt.Sprintf("You lost! Time: %s", window.game.FormattedTime())

	default:
		return fmt.Sprintf("Flags: %d  Time: %s", window.game.FlagsRemaining(), window.game.FormattedTime())
	}
}

func (window *BoardWindow) handleKey(key string) tea.Cmd {
	game := window.game

	switch key {
	case "up", "k":
		if window.CursorRow == 0 {
			window.CursorRow = game.Rows() - 1
			break
		}
		window.CursorRow--

	case "down", "j":
		if window.CursorRow == game.Rows()-1 {
			window.CursorRow = 0
			break
		}
		window.CursorRow++

	case "left", "h":
		if window.CursorCol == 0 {
			window.CursorCol = game.Cols() - 1
			break
		}
		window.CursorCol--

	case "right", "l":
		if window.CursorCol == game.Cols()-1 {
			window.CursorCol = 0
			break
		}
		window.CursorCol++

	case "space", "enter":
		game.Reveal(window.CursorRow, window.CursorCol)

	case "f":
		game.ToggleFlag(window.CursorRow, window.CursorCol)

	case "r":
		window.NewGame(window.settings)

	case "q":
		return func() tea.Msg { return exitBoardMsg{} }
	}

	return nil
}

func (window *BoardWindow) startClock() tea.Cmd {
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
