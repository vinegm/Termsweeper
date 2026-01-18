package src

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// MenuWindow implements the menu screen and holds its view state.
type MenuWindow struct {
	Choice int

	minWidth  int
	minHeight int
}

func (window *MenuWindow) Render(model *model) string {
	var sb strings.Builder

	sb.WriteString("Termsweeper\n")

	items := []string{"Start Game", "Quit"}
	for i, it := range items {
		sb.WriteString("\n")
		if window.Choice == i {
			sb.WriteString(selectedMenuStyle.Render(it))
			continue
		}

		sb.WriteString(it)
	}

	return sb.String()
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
			if window.Choice < 1 {
				window.Choice++
			}

		case "enter", " ":
			switch window.Choice {
			case 0:
				model.inGame = true
				model.game.minesPlaced = false
				model.CurrentWindow = model.BoardWin

			case 1:
				return tea.Quit
			}

		case "q":
			return tea.Quit
		}
	}

	return nil
}

func (window *MenuWindow) MinSize() (int, int) { return window.minWidth, window.minHeight }
