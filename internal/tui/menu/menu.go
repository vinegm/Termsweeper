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
	boardRows
	boardCols
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
	// Asks for control to start a new game with these settings.
	StartGameMsg struct{ Settings minesweeper.Settings }

	// Asks for control to return to the game already in progress.
	ContinueGameMsg struct{}
)

// Window is the main menu window for the game, allowing
// the user to start a new game or change settings.
type Window struct {
	Choice      int
	OngoingGame bool
	EditingSeed bool
	Settings    minesweeper.Settings

	theme *theme.Theme
}

// New creates a new Window with the provided theme and settings.
func New(th *theme.Theme, settings minesweeper.Settings) *Window {
	return &Window{
		Settings: settings,
		theme:    th,
	}
}

func (window *Window) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if window.EditingSeed {
			window.editSeed(msg)
			return nil
		}

		return window.handleKey(strings.ToLower(msg.String()))
	}

	return nil
}

// View renders the menu window as a string.
func (window *Window) View(_, _ int) string {
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
		menuItem{fmt.Sprintf("Board Rows: %d", window.Settings.Rows), boardRows},
		menuItem{fmt.Sprintf("Board Cols: %d", window.Settings.Cols), boardCols},
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

	switch action {
	case boardRows:
		settings.Rows = clampDim(settings.Rows + step)

	case boardCols:
		settings.Cols = clampDim(settings.Cols + step)

	case mineCount:
		settings.Mines += step
	}

	settings.Mines = min(max(settings.Mines, minesweeper.MinMines), settings.MaxMines())
}

// clampDim keeps a board dimension within its limits.
func clampDim(dim int) int {
	return min(max(dim, minesweeper.MinDim), minesweeper.MaxDim)
}
