package menu

import (
	"fmt"
	"strings"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"

	"github.com/vinegm/termsweeper/internal/minesweeper"
	"github.com/vinegm/termsweeper/internal/tui/theme"
)

type menuOption int

// Menu options that can be selected by the user.
const (
	continueGame menuOption = iota
	startNewGame
	boardSize
	mineCount
	seed
	quit
)

// Longest seed that can be typed in the menu.
const maxSeedLen = 32

// A single selectable row of the menu.
type menuItem struct {
	label  string
	action menuOption
}

type (
	// StartGameMsg asks for control to start a new game with these settings.
	StartGameMsg struct{ Settings minesweeper.Settings }

	// ContinueGameMsg asks for control to return to the game already in progress.
	ContinueGameMsg struct{}
)

// Window is the main menu window for the game, allowing
// the user to start a new game or change settings.
type Window struct {
	Choice      int
	OngoingGame bool
	EditingSeed bool
	EditingSize bool

	Settings   minesweeper.Settings
	sizeBackup minesweeper.Settings

	theme *theme.Theme
}

// New creates a new Window with the provided theme and settings.
func New(th *theme.Theme, settings minesweeper.Settings) *Window {
	return &Window{
		Settings: settings,
		theme:    th,
	}
}

// Update handles the keys, typed into the seed while editing it.
func (window *Window) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if window.EditingSeed {
			window.editSeed(msg)
			return nil
		}

		if window.EditingSize {
			window.editSize(strings.ToLower(msg.String()))
			return nil
		}

		return window.handleKey(strings.ToLower(msg.String()))
	}

	return nil
}

// SetSize does nothing, the menu is drawn at its own size.
func (window *Window) SetSize(_, _ int) {}

// View renders the menu window as a string.
func (window *Window) View() string {
	var sb strings.Builder

	sb.WriteString(window.theme.TitleStyle.Render(getMenuTitle()))
	sb.WriteString("\n")

	for i, item := range window.items() {
		sb.WriteString("\n")
		if window.Choice == i {
			sb.WriteString(window.theme.SelectedMenuTextStyle.Render(item.label))
			continue
		}

		sb.WriteString(window.theme.TextStyle.Render(item.label))
	}

	return window.theme.WindowStyle.Render(sb.String())
}

func getMenuTitle() string {
	art := `
▗▄▄▄▖▗▞▀▚▖ ▄▄▄ ▄▄▄▄   ▗▄▄▖▄   ▄ ▗▞▀▚▖▗▞▀▚▖▄▄▄▄  ▗▞▀▚▖ ▄▄▄
  █  ▐▛▀▀▘█    █ █ █ ▐▌   █ ▄ █ ▐▛▀▀▘▐▛▀▀▘█   █ ▐▛▀▀▘█   
  █  ▝▚▄▄▖█    █   █  ▝▀▚▖█▄█▄█ ▝▚▄▄▖▝▚▄▄▖█▄▄▄▀ ▝▚▄▄▖█   
  █                  ▗▄▄▞▘                █              
                                          ▀              `

	return art[1:] // Remove leading newline
}

func (window *Window) items() []menuItem {
	items := []menuItem{}

	if window.OngoingGame {
		items = append(items,
			menuItem{"Continue", continueGame},
			menuItem{"New Game", startNewGame},
		)
	} else {
		items = append(items, menuItem{"Start Game", startNewGame})
	}

	return append(items,
		menuItem{window.sizeLabel(), boardSize},
		menuItem{fmt.Sprintf("Mines: %d", window.Settings.Mines), mineCount},
		menuItem{window.seedLabel(), seed},
		menuItem{"Quit", quit},
	)
}

func (window *Window) handleKey(key string) tea.Cmd {
	items := window.items()

	switch key {
	case "up", "k":
		if window.Choice > 0 {
			window.Choice--
		}

	case "down", "j":
		if window.Choice < len(items)-1 {
			window.Choice++
		}

	case "left":
		window.adjust(items[window.Choice].action, -1)

	case "right":
		window.adjust(items[window.Choice].action, 1)

	case "enter", "space":
		switch items[window.Choice].action {
		case continueGame:
			return func() tea.Msg {
				return ContinueGameMsg{}
			}

		case startNewGame:
			chosen := window.Settings
			return func() tea.Msg {
				return StartGameMsg{Settings: chosen}
			}

		case boardSize:
			window.sizeBackup = window.Settings
			window.EditingSize = true

		case seed:
			window.EditingSeed = true

		case quit:
			return tea.Quit
		}

	case "q":
		return tea.Quit
	}

	return nil
}

// editSeed handles key presses while editing the seed.
func (window *Window) editSeed(msg tea.KeyPressMsg) {
	text := &window.Settings.Seed

	switch msg.String() {
	case "enter", "esc":
		window.EditingSeed = false

	case "backspace":
		_, size := utf8.DecodeLastRuneInString(*text)
		*text = (*text)[:len(*text)-size]

	default:
		isText := msg.Text != ""
		fits := utf8.RuneCountInString(*text+msg.Text) <= maxSeedLen

		if isText && fits {
			*text += msg.Text
		}
	}
}

// editSize handles key presses while resizing the board.
func (window *Window) editSize(key string) {
	settings := &window.Settings

	switch key {
	case "left", "h":
		settings.Cols = clampDim(settings.Cols - 1)

	case "right", "l":
		settings.Cols = clampDim(settings.Cols + 1)

	case "up", "k":
		settings.Rows = clampDim(settings.Rows - 1)

	case "down", "j":
		settings.Rows = clampDim(settings.Rows + 1)

	case "enter", "space":
		window.EditingSize = false

	case "esc":
		*settings = window.sizeBackup
		window.EditingSize = false
	}

	settings.Mines = min(settings.Mines, settings.MaxMines())
}

func (window *Window) sizeLabel() string {
	size := fmt.Sprintf("%dx%d", window.Settings.Cols, window.Settings.Rows)

	if window.EditingSize {
		return "Board Size: ◂ " + size + " ▸"
	}

	return "Board Size: " + size
}

func (window *Window) seedLabel() string {
	if window.EditingSeed {
		return "Seed: " + window.Settings.Seed + "_"
	}

	if window.Settings.Seed == "" {
		return "Seed: random"
	}

	return "Seed: " + window.Settings.Seed
}

// adjust steps the setting the given menu row edits.
func (window *Window) adjust(action menuOption, step int) {
	settings := &window.Settings

	if action == mineCount {
		settings.Mines = min(max(settings.Mines+step, minesweeper.MinMines), settings.MaxMines())
	}
}

// clampDim keeps a board dimension within its limits.
func clampDim(dim int) int {
	return min(max(dim, minesweeper.MinDim), minesweeper.MaxDim)
}
