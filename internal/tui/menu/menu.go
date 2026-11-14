package menu

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"charm.land/bubbles/v2/key"
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
			window.editSize(lowered(strings.ToLower(msg.String())))
			return nil
		}

		return window.handleKey(lowered(strings.ToLower(msg.String())))
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
	sb.WriteString(window.theme.MutedStyle.Render("Minesweeper"))
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
		menuItem{fmt.Sprintf("Mines: ◂ %d ▸", window.Settings.Mines), mineCount},
		menuItem{window.seedLabel(), seed},
		menuItem{"Quit", quit},
	)
}

func (window *Window) handleKey(pressed lowered) tea.Cmd {
	items := window.items()

	switch {
	case key.Matches(pressed, keys.Up):
		if window.Choice > 0 {
			window.Choice--
		}

	case key.Matches(pressed, keys.Down):
		if window.Choice < len(items)-1 {
			window.Choice++
		}

	case key.Matches(pressed, keys.Left):
		window.adjust(items[window.Choice].action, -1)

	case key.Matches(pressed, keys.Right):
		window.adjust(items[window.Choice].action, 1)

	case key.Matches(pressed, keys.Select):
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

	case key.Matches(pressed, keys.Quit):
		return tea.Quit
	}

	return nil
}

// editSeed handles key presses while editing the seed.
func (window *Window) editSeed(msg tea.KeyPressMsg) {
	text := &window.Settings.Seed

	switch {
	case key.Matches(msg, keys.SeedDone):
		window.EditingSeed = false

	case key.Matches(msg, keys.Delete):
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
func (window *Window) editSize(pressed lowered) {
	settings := &window.Settings

	switch {
	case key.Matches(pressed, keys.Left):
		settings.Cols = clampDim(settings.Cols - 1)

	case key.Matches(pressed, keys.Right):
		settings.Cols = clampDim(settings.Cols + 1)

	case key.Matches(pressed, keys.Up):
		settings.Rows = clampDim(settings.Rows - 1)

	case key.Matches(pressed, keys.Down):
		settings.Rows = clampDim(settings.Rows + 1)

	case key.Matches(pressed, keys.SizeDone):
		window.EditingSize = false

	case key.Matches(pressed, keys.Cancel):
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

// ShortHelp is the keys for the key hints bar, for the selected row.
func (window *Window) ShortHelp() []key.Binding {
	if window.EditingSeed {
		return []key.Binding{keys.SeedDone, keys.Delete}
	}

	if window.EditingSize {
		return []key.Binding{keys.Resize, keys.SizeDone, keys.Cancel}
	}

	bindings := []key.Binding{keys.Move}

	switch window.items()[window.Choice].action {
	case mineCount:
		bindings = append(bindings, keys.Change)

	case boardSize, seed:
		edit := keys.Select
		edit.SetHelp("enter", "Edit")
		bindings = append(bindings, edit)

	default:
		bindings = append(bindings, keys.Select)
	}

	return append(bindings, keys.Quit)
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
