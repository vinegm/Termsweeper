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
//	WARN: Must be able to fit any window.
const (
	minWidth  = 61
	minHeight = 13
)

// Window is the interface that holds methods to handle
// updates and rendering for different windows.
type Window interface {
	Update(msg tea.Msg) tea.Cmd
	View(width, height int) string
}

// Model represents the state of the application.
type Model struct {
	width  int
	height int

	currentWindow Window
	menu          *menu.Window
	board         *board.Window
}

// New creates a new model with initial state.
func New(cfg config.Config, settings minesweeper.Settings) Model {
	th := theme.NewTheme(cfg)
	menuWindow := menu.New(th, settings)

	return Model{
		currentWindow: menuWindow,
		menu:          menuWindow,
		board:         board.New(th, boardOptions(cfg)),
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

// Tea update loop
func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

// Tea view rendering, handles terminal size checks
// and delegates to the current window.
func (model Model) View() tea.View {
	content := model.viewContent()

	view := tea.NewView(
		lg.Place(
			model.width,
			model.height,
			lg.Center,
			lg.Center,
			content,
		),
	)
	view.AltScreen = true

	return view
}

func (model Model) viewContent() string {
	if model.width < minWidth || model.height < minHeight {
		return fmt.Sprintf("Terminal too small — need at least %dx%d", minWidth, minHeight)
	}

	// Render centered base panel
	return model.currentWindow.View(model.width, model.height)
}
