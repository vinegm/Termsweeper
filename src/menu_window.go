package src

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Implements the menu screen and holds its view state.
type MenuWindow struct {
	Choice int

	minWidth  int
	minHeight int
}

// Initializes a new MenuWindow with the first menu item selected and minimum size set.
func NewMenuWindow() *MenuWindow {
	return &MenuWindow{Choice: 0, minWidth: 30, minHeight: 10}
}

func (window *MenuWindow) Render(model *model) string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render(getMenuTitle()) + "\n")

	rows := AppConfig.BoardRows
	cols := AppConfig.BoardCols
	mines := AppConfig.MineCount

	items := []string{
		"Start Game",
		"Board Rows: " + fmt.Sprintf("%d", rows),
		"Board Cols: " + fmt.Sprintf("%d", cols),
		"Mines: " + fmt.Sprintf("%d", mines),
		"Quit",
	}

	for i, it := range items {
		sb.WriteString("\n")
		if window.Choice == i {
			sb.WriteString(selectedMenuTextStyle.Render(it))
			continue
		}

		sb.WriteString(textStyle.Render(it))
	}

	return windowStyle.Render(sb.String())
}

func (window *MenuWindow) HandleInput(model *model, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if window.Choice > 0 {
				window.Choice--
			}

		case "down", "j":
			if window.Choice < 4 {
				window.Choice++
			}

		case "left":
			switch window.Choice {
			case 1:
				if AppConfig.BoardRows > 5 {
					AppConfig.BoardRows--
					max := AppConfig.BoardRows * AppConfig.BoardCols
					if AppConfig.MineCount >= max {
						AppConfig.MineCount = max - 1
					}
				}

			case 2:
				if AppConfig.BoardCols > 5 {
					AppConfig.BoardCols--
					max := AppConfig.BoardRows * AppConfig.BoardCols
					if AppConfig.MineCount >= max {
						AppConfig.MineCount = max - 1
					}
				}

			case 3:
				if AppConfig.MineCount > 1 {
					AppConfig.MineCount--
				}
			}

		case "right":
			switch window.Choice {
			case 1:
				if AppConfig.BoardRows < 40 {
					AppConfig.BoardRows++
				}

			case 2:
				if AppConfig.BoardCols < 40 {
					AppConfig.BoardCols++
				}

			case 3:
				max := AppConfig.BoardRows * AppConfig.BoardCols
				if AppConfig.MineCount < max-1 {
					AppConfig.MineCount++
				}
			}

		case "enter", " ":
			switch window.Choice {
			case 0:
				model.game = newGame()
				model.CurrentWindow = model.BoardWin

			case 4:
				return tea.Quit
			}

		case "q":
			return tea.Quit
		}
	}

	return nil
}

func (window *MenuWindow) MinSize(_ *model) (int, int) { return window.minWidth, window.minHeight }

func getMenuTitle() string {
	art := `
▗▄▄▄▖▗▞▀▚▖ ▄▄▄ ▄▄▄▄   ▗▄▄▖▄   ▄ ▗▞▀▚▖▗▞▀▚▖▄▄▄▄  ▗▞▀▚▖ ▄▄▄
  █  ▐▛▀▀▘█    █ █ █ ▐▌   █ ▄ █ ▐▛▀▀▘▐▛▀▀▘█   █ ▐▛▀▀▘█   
  █  ▝▚▄▄▖█    █   █  ▝▀▚▖█▄█▄█ ▝▚▄▄▖▝▚▄▄▖█▄▄▄▀ ▝▚▄▄▖█   
  █                  ▗▄▄▞▘                █              
                                          ▀              `

	return art[1:] // Remove leading newline
}
