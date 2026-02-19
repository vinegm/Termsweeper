package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type menuOption int

const (
	startGame menuOption = iota
	boardRows
	boardCols
	mineCount
	quit
)

// Implements the menu screen and holds its view state.
type MenuWindow struct {
	Choice menuOption

	minWidth  int
	minHeight int
}

// Initializes a new MenuWindow with the first menu item selected and minimum size set.
func NewMenuWindow() *MenuWindow {
	return &MenuWindow{Choice: startGame, minWidth: 30, minHeight: 10}
}

// The menu lines, in the order the menuOption constants declare them.
func menuItems(model *Model) []string {
	rows := model.settings.Rows
	cols := model.settings.Cols
	mines := model.settings.Mines

	return []string{
		"Start Game",
		"Board Rows: " + fmt.Sprintf("%d", rows),
		"Board Cols: " + fmt.Sprintf("%d", cols),
		"Mines: " + fmt.Sprintf("%d", mines),
		"Quit",
	}
}

func (window *MenuWindow) Render(model *Model) string {
	var sb strings.Builder

	sb.WriteString(model.theme.titleStyle.Render(getMenuTitle()))
	sb.WriteString("\n")

	for i, it := range menuItems(model) {
		sb.WriteString("\n")
		if window.Choice == menuOption(i) {
			sb.WriteString(model.theme.selectedMenuTextStyle.Render(it))
			continue
		}

		sb.WriteString(model.theme.textStyle.Render(it))
	}

	return model.theme.windowStyle.Render(sb.String())
}

func (window *MenuWindow) HandleInput(model *Model, msg string) tea.Cmd {
	settings := &model.settings

	switch msg {
	case "up", "k":
		if window.Choice > 0 {
			window.Choice--
		}

	case "down", "j":
		if window.Choice < menuOption(len(menuItems(model))-1) {
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
		maxRows, maxCols := maxBoardSize(model.width, model.height)

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
			model.newGame()
			model.CurrentWindow = model.BoardWin

		case quit:
			return tea.Quit
		}

	case "q":
		return tea.Quit
	}

	return nil
}

func (window *MenuWindow) MinSize(_ *Model) (int, int) { return window.minWidth, window.minHeight }

func getMenuTitle() string {
	art := `
▗▄▄▄▖▗▞▀▚▖ ▄▄▄ ▄▄▄▄   ▗▄▄▖▄   ▄ ▗▞▀▚▖▗▞▀▚▖▄▄▄▄  ▗▞▀▚▖ ▄▄▄
  █  ▐▛▀▀▘█    █ █ █ ▐▌   █ ▄ █ ▐▛▀▀▘▐▛▀▀▘█   █ ▐▛▀▀▘█   
  █  ▝▚▄▄▖█    █   █  ▝▀▚▖█▄█▄█ ▝▚▄▄▖▝▚▄▄▖█▄▄▄▀ ▝▚▄▄▖█   
  █                  ▗▄▄▞▘                █              
                                          ▀              `

	return art[1:] // Remove leading newline
}
