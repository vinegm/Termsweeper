package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/termsweeper/internal/tui/board"
	"github.com/vinegm/termsweeper/internal/tui/menu"

	"github.com/vinegm/termsweeper/internal/config"
	"github.com/vinegm/termsweeper/internal/minesweeper"
	"github.com/vinegm/termsweeper/internal/tui/theme"
)

// Smallest terminal the app renders correctly in.
//
//	WARN: Must be able to fit any window, and the key hints bar under it.
const (
	minWidth  = 61
	minHeight = 15 + hintsHeight
)

// Rows the key hints bar takes at the bottom of the screen.
const hintsHeight = 1

// Window is a screen of the app. SetSize is called after
// every Update, before View draws at that size.
type Window interface {
	Update(msg tea.Msg) tea.Cmd
	SetSize(width, height int)
	View() string
	Hints() []theme.Hint
}

// Model represents the state of the application.
type Model struct {
	width  int
	height int

	currentWindow Window
	menu          *menu.Window
	board         *board.Window

	theme *theme.Theme
}

// New creates a new model with initial state.
func New(cfg config.Config, settings minesweeper.Settings) Model {
	th := theme.NewTheme(cfg)
	menuWindow := menu.New(th, settings)

	return Model{
		currentWindow: menuWindow,
		menu:          menuWindow,
		board:         board.New(th, boardOptions(cfg)),
		theme:         th,
	}
}

func boardOptions(cfg config.Config) board.Options {
	return board.Options{
		Highlight:             cfg.Layout.Highlight,
		Numbers:               cfg.Layout.Numbers,
		RelativeLineNumbers:   cfg.Layout.RelativeLineNumbers,
		RelativeColumnNumbers: cfg.Layout.RelativeColumnNumbers,
		Scrolloff:             cfg.Behaviour.Scrolloff,
		SideScrolloff:         cfg.Behaviour.SideScrolloff,
	}
}

// Init initializes the model and returns an initial command.
//
//	WARN: This function is required by the Bubble Tea framework.
func (model Model) Init() tea.Cmd {
	return nil
}

// Update handles a message, then sizes the current window for View.
func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := model.update(msg)
	model.currentWindow.SetSize(model.width, model.height-hintsHeight)

	return model, cmd
}

func (model Model) update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		return model, nil

	case menu.StartGameMsg:
		model.board.NewGame(msg.Settings)
		model.menu.OngoingGame = true
		model.currentWindow = model.board
		return model, nil

	case menu.ContinueGameMsg:
		model.currentWindow = model.board
		return model, nil

	case board.ExitMsg:
		model.currentWindow = model.menu
		return model, nil
	}

	return model, model.currentWindow.Update(msg)
}

// View draws the current window centred on the screen, with the
// key hints bar under it, or a notice if the terminal is too small.
func (model Model) View() tea.View {
	view := tea.NewView(model.viewContent())
	view.AltScreen = true

	return view
}

func (model Model) viewContent() string {
	if model.width < minWidth || model.height < minHeight {
		return lg.Place(
			model.width,
			model.height,
			lg.Center,
			lg.Center,
			fmt.Sprintf("Terminal too small — need at least %dx%d", minWidth, minHeight),
		)
	}

	// Render centered base panel
	window := lg.Place(
		model.width,
		model.height-hintsHeight,
		lg.Center,
		lg.Center,
		model.currentWindow.View(),
	)

	hints := model.theme.RenderHints(model.currentWindow.Hints(), model.width)

	return lg.JoinVertical(lg.Left, window, hints)
}
