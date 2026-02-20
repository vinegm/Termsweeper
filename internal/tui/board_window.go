package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"termsweeper/internal/minesweeper"
)

// Space the board window needs on top of the board itself.
const (
	boardChromeWidth  = 4 // borders + padding
	boardChromeHeight = 6 // title + hint + borders + padding
)

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

	theme    *Theme
	game     *minesweeper.Game
	settings minesweeper.Settings

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
	window.settings = settings
	window.game = minesweeper.New(settings)

	window.clockTag++
	window.ticking = false
}

func (window *BoardWindow) Update(msg tea.Msg) (Window, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if cmd := window.handleKey(strings.ToLower(msg.String())); cmd != nil {
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
	var sb strings.Builder

	sb.WriteString(window.theme.titleStyle.Render("Termsweeper"))
	sb.WriteString("\n")
	sb.WriteString(window.theme.textStyle.Render(window.hint()))
	sb.WriteString("\n\n")

	for row := range window.game.Rows() {
		for col := range window.game.Cols() {
			cell := window.game.Cell(row, col)
			char := window.theme.cellChar(cell)

			isCursorRow := window.CursorRow == row
			isCursorCol := window.CursorCol == col
			isSelectedCell := isCursorRow && isCursorCol
			style := window.theme.cellStyle(cell, isSelectedCell)

			sb.WriteString(style.Render(char + " "))
		}
		sb.WriteString("\n")
	}

	return window.theme.windowStyle.Render(sb.String())
}

func (window *BoardWindow) MinSize() (int, int) {
	game := window.game

	boardWidth := game.Cols() * 2 // char + space
	minWidth := boardWidth + boardChromeWidth
	minHeight := game.Rows() + boardChromeHeight

	return minWidth, minHeight
}

func maxBoardSize(width, height int) (rows, cols int) {
	return height - boardChromeHeight, (width - boardChromeWidth) / 2
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
