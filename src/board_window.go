package src

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// BoardWindow implements the game/board screen and holds its view state.
type BoardWindow struct {
	CursorRow int
	CursorCol int

	minWidth  int
	minHeight int
}

func (window *BoardWindow) Render(model *model) string {
	var sb strings.Builder

	sb.WriteString("Termsweeper\n")
	sb.WriteString(getHint(model) + "\n\n")

	for row := range model.game.board {
		for col := range model.game.board[row] {
			cell := model.game.board[row][col]
			char := cell.char()

			isCursorRow := window.CursorRow == row
			isCursorCol := window.CursorCol == col
			isSelectedCell := isCursorRow && isCursorCol
			style := cellStyle(row, col, isSelectedCell)

			sb.WriteString(style.Render(char + " "))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (window *BoardWindow) HandleInput(model *model, msg tea.Msg) tea.Cmd {
	game := &model.game

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if window.CursorRow > 0 {
				window.CursorRow--
			}

		case "down", "j":
			if window.CursorRow < rows-1 {
				window.CursorRow++
			}

		case "left", "h":
			if window.CursorCol > 0 {
				window.CursorCol--
			}

		case "right", "l":
			if window.CursorCol < cols-1 {
				window.CursorCol++
			}

		case " ", "enter":
			if !game.minesPlaced {
				game.placeMines(window.CursorRow, window.CursorCol)
				game.minesPlaced = true
			}

			game.reveal(window.CursorRow, window.CursorCol)
			if game.state == playing {
				game.checkWin()
			}

		case "f":
			game.toggleFlag(window.CursorRow, window.CursorCol)

		case "r":
			model.inGame = true
			model.game = InitialModel().game
			model.CurrentWindow = model.BoardWin

		case "q":
			model.inGame = false
			model.CurrentWindow = model.MenuWin
		}
	}
	return nil
}

func (window *BoardWindow) MinSize() (int, int) { return window.minWidth, window.minHeight }

func getHint(model *model) string {
	switch model.game.state {
	case won:
		return "You won! Press q to quit or r to restart."

	case lost:
		return "You lost! Press q to quit or r to restart."

	default:
		return fmt.Sprintf("Flags: %d", model.game.flagsRemaining())
	}
}
