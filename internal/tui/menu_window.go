package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"termsweeper/internal/minesweeper"
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
type startGameMsg struct{ settings minesweeper.Settings }

// MenuWindow is the main menu window for the game, allowing
// the user to start a new game or change settings.
type MenuWindow struct {
	Choice   menuOption
	Settings minesweeper.Settings

	theme         *Theme
	width, height int
	minWidth      int
	minHeight     int
}

// NewMenuWindow creates a new MenuWindow with the provided theme and settings.
func NewMenuWindow(theme *Theme, settings minesweeper.Settings) *MenuWindow {
	return &MenuWindow{
		Choice:    startGame,
		Settings:  settings,
		theme:     theme,
		minWidth:  30,
		minHeight: 10,
	}
}

func (window *MenuWindow) Update(msg tea.Msg) (Window, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		window.width = msg.Width
		window.height = msg.Height

	case tea.KeyPressMsg:
		return window, window.handleKey(strings.ToLower(msg.String()))
	}

	return window, nil
}

// View renders the menu window as a string.
func (window *MenuWindow) View() string {
	var sb strings.Builder

	sb.WriteString(window.theme.titleStyle.Render(getMenuTitle()))
	sb.WriteString("\n")

	for i, it := range window.items() {
		sb.WriteString("\n")
		if window.Choice == menuOption(i) {
			sb.WriteString(window.theme.selectedMenuTextStyle.Render(it))
			continue
		}

		sb.WriteString(window.theme.textStyle.Render(it))
	}

	return window.theme.windowStyle.Render(sb.String())
}

// MinSize returns the minimum width and height of the menu window.
func (window *MenuWindow) MinSize() (int, int) {
	return window.minWidth, window.minHeight
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

func (window *MenuWindow) items() []string {
	return []string{
		"Start Game",
		"Board Rows: " + fmt.Sprintf("%d", window.Settings.Rows),
		"Board Cols: " + fmt.Sprintf("%d", window.Settings.Cols),
		"Mines: " + fmt.Sprintf("%d", window.Settings.Mines),
		"Quit",
	}
}

func (window *MenuWindow) handleKey(key string) tea.Cmd {
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
		maxRows, maxCols := maxBoardSize(window.width, window.height)

		switch window.Choice {
		case boardRows:
			if settings.Rows < maxRows {
				settings.Rows++
			}

		case boardCols:
			if settings.Cols < maxCols {
				settings.Cols++
			}

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
				return startGameMsg{settings: chosen}
			}

		case quit:
			return tea.Quit
		}

	case "q":
		return tea.Quit
	}

	return nil
}
