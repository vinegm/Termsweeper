package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

// Tea update loop
func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		return model, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return model, tea.Quit

		case "up", "k":
			if !model.inGame {
				if model.menuChoice > 0 {
					model.menuChoice--
				}
				return model, nil
			}

			if model.game.cursorRow > 0 {
				model.game.cursorRow--
			}

		case "down", "j":
			if !model.inGame {
				if model.menuChoice < 1 {
					model.menuChoice++
				}
				return model, nil
			}

			if model.game.cursorRow < rows-1 {
				model.game.cursorRow++
			}

		case "left", "h":
			if model.game.cursorCol > 0 {
				model.game.cursorCol--
			}

		case "right", "l":
			if model.game.cursorCol < cols-1 {
				model.game.cursorCol++
			}

		case " ", "enter":
			if !model.inGame {
				switch model.menuChoice {
				case 0:
					model.inGame = true
					model.game.minesPlaced = false
				case 1:
					return model, tea.Quit
				}

				return model, nil
			}

			if model.game.state != playing {
				return model, nil
			}

			if !model.game.minesPlaced {
				model.game.placeMines(model.game.cursorRow, model.game.cursorCol)
				model.game.minesPlaced = true
			}

			model.game.reveal(model.game.cursorRow, model.game.cursorCol)
			if model.game.state == playing {
				model.game.checkWin()
			}

		case "f":
			if model.game.state != playing {
				return model, nil
			}

			model.game.toggleFlag(model.game.cursorRow, model.game.cursorCol)
			model.game.checkWin()
		}
	}

	return model, nil
}

// Tea view rendering, handles terminal size checks and delegates to menu or game renderers.
func (model model) View() string {
	var minW, minH int
	if model.inGame {
		minW, minH = model.boardMinW, model.boardMinH
	} else {
		minW, minH = model.menuMinW, model.menuMinH
	}

	if model.width > 0 && model.height > 0 && (model.width < minW || model.height < minH) {
		warn := fmt.Sprintf("Terminal too small — need at least %dx%d", minW, minH)
		return lg.Place(model.width, model.height, lg.Center, lg.Center, warn)
	}

	var panel string
	if model.inGame {
		panel = renderGame(model)
	} else {
		panel = renderMenu(model)
	}

	if model.width > 0 && model.height > 0 {
		return lg.Place(model.width, model.height, lg.Center, lg.Center, panel)
	}

	return panel
}
