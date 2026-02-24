package board

import (
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
	window.settings = settings
	window.countBuffer = ""

	window.CursorRow = 0
	window.CursorCol = 0

	window.offsetRow = 0
	window.offsetCol = 0

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
