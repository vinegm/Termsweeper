package menu

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"termsweeper/internal/minesweeper"
	"termsweeper/internal/tui/theme"
)

type menuOption int

const (
	startGame menuOption = iota
	boardRows
	boardCols
	mineCount
	quit
)

// Asks for control to start the game with these settings.
type StartGameMsg struct{ Settings minesweeper.Settings }

// Window is the main menu window for the game, allowing
// the user to start a new game or change settings.
type Window struct {
	Choice   menuOption
	Settings minesweeper.Settings

	theme *theme.Theme
}

// New creates a new Window with the provided theme and settings.
func New(th *theme.Theme, settings minesweeper.Settings) *Window {
	return &Window{
		Choice:   startGame,
		Settings: settings,
		theme:    th,
	}
}

func (window *Window) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return window.handleKey(strings.ToLower(msg.String()))
	}

	return nil
}

// View renders the menu window as a string.
func (window *Window) View(_, _ int) string {
	var sb strings.Builder

	sb.WriteString(window.theme.TitleStyle.Render(getMenuTitle()))
	sb.WriteString("\n")

	for i, it := range window.items() {
		sb.WriteString("\n")
		if window.Choice == menuOption(i) {
			sb.WriteString(window.theme.SelectedMenuTextStyle.Render(it))
			continue
		}

		sb.WriteString(window.theme.TextStyle.Render(it))
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

func (window *Window) items() []string {
	return []string{
		"Start Game",
		"Board Rows: " + fmt.Sprintf("%d", window.Settings.Rows),
		"Board Cols: " + fmt.Sprintf("%d", window.Settings.Cols),
		"Mines: " + fmt.Sprintf("%d", window.Settings.Mines),
		"Quit",
	}
}

func (window *Window) handleKey(key string) tea.Cmd {
	settings := &window.Settings

	switch key {
	case "up", "k":
		if window.Choice > 0 {
			window.Choice--
		}

	case "down", "j":
		if window.Choice < menuOption(len(window.items())-1) {
			window.Choice++
		}

	case "left":
		switch window.Choice {
		case boardRows:
			if settings.Rows > 5 {
				settings.Rows--
				max := settings.Rows * settings.Cols

				if settings.Mines >= max {
					settings.Mines = max - 1
				}
			}

		case boardCols:
			if settings.Cols > 5 {
				settings.Cols--
				max := settings.Rows * settings.Cols
				if settings.Mines >= max {
					settings.Mines = max - 1
				}
			}

		case mineCount:
			if settings.Mines > 1 {
				settings.Mines--
			}
		}

	case "right":
		switch window.Choice {
		case boardRows:
			settings.Rows++

		case boardCols:
			settings.Cols++

		case mineCount:
			max := settings.Rows * settings.Cols
			if settings.Mines < max-1 {
				settings.Mines++
			}
		}

	case "enter", "space":
		switch window.Choice {
		case startGame:
			chosen := window.Settings
			return func() tea.Msg {
				return StartGameMsg{Settings: chosen}
			}

		case quit:
			return tea.Quit
		}

	case "q":
		return tea.Quit
	}

	return nil
}
