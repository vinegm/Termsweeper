package src

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Implements the game/board screen and holds its view state.
type BoardWindow struct {
	CursorRow int
	CursorCol int
}

// Initializes a new BoardWindow with the cursor at the top-left corner.
func NewBoardWindow() *BoardWindow {
	return &BoardWindow{CursorRow: 0, CursorCol: 0}
}

func (window *BoardWindow) Render(model *model) string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("Termsweeper") + "\n")
	sb.WriteString(textStyle.Render(getHint(model)) + "\n\n")

	for row := range model.game.board {
		for col := range model.game.board[row] {
			cell := model.game.board[row][col]
			char := cell.char()

			isCursorRow := window.CursorRow == row
			isCursorCol := window.CursorCol == col
			isSelectedCell := isCursorRow && isCursorCol
			style := cellStyle(cell, isSelectedCell)

			sb.WriteString(style.Render(char + " "))
		}
		sb.WriteString("\n")
	}

	return windowStyle.Render(sb.String())
}

func (window *BoardWindow) HandleInput(model *model, msg string) tea.Cmd {
	game := model.game

	switch msg {
	case "up", "k":
		if window.CursorRow == 0 {
			window.CursorRow = game.rows-1
			break
		}
		window.CursorRow--

	case "down", "j":
		if window.CursorRow == game.rows-1 {
			window.CursorRow = 0
			break
		}
		window.CursorRow++

	case "left", "h":
		if window.CursorCol == 0 {
			window.CursorCol = game.cols-1
			break
		}
		window.CursorCol--

	case "right", "l":
		if window.CursorCol == game.cols-1 {
			window.CursorCol = 0
			break
		}
		window.CursorCol++

	case " ", "enter":
		game.reveal(window.CursorRow, window.CursorCol)

	case "f":
		game.toggleFlag(window.CursorRow, window.CursorCol)

	case "r":
		model.game = newGame()

	case "q":
		window.CursorCol = 0
		window.CursorRow = 0

		model.game = newGame()
		model.CurrentWindow = model.MenuWin
	}

	return nil
}

func (window *BoardWindow) MinSize(model *model) (int, int) {
	game := model.game

	boardWidth := game.cols * 2 // char + space
	minWidth := boardWidth + 4  // board + borders + padding
	minHeight := game.rows + 6  // title + hint + borders + padding

	return minWidth, minHeight
}

// Returns a hint string based on the current game state.
func getHint(model *model) string {
	switch model.game.state {
	case won:
		return fmt.Sprintf("You won! Time: %s", model.game.FormattedTime())

	case lost:
		return fmt.Sprintf("You lost! Time: %s", model.game.FormattedTime())

	default:
		return fmt.Sprintf("Flags: %d  Time: %s", model.game.flagsRemaining(), model.game.FormattedTime())
	}
}
