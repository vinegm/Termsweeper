package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"termsweeper/internal/minesweeper"
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

func (window *BoardWindow) Render(model *Model) string {
	var sb strings.Builder

	sb.WriteString(model.theme.titleStyle.Render("Termsweeper") + "\n")
	sb.WriteString(model.theme.textStyle.Render(getHint(model)) + "\n\n")

	for row := range model.game.Rows() {
		for col := range model.game.Cols() {
			cell := model.game.Cell(row, col)
			char := model.theme.cellChar(cell)

			isCursorRow := window.CursorRow == row
			isCursorCol := window.CursorCol == col
			isSelectedCell := isCursorRow && isCursorCol
			style := model.theme.cellStyle(cell, isSelectedCell)

			sb.WriteString(style.Render(char + " "))
		}
		sb.WriteString("\n")
	}

	return model.theme.windowStyle.Render(sb.String())
}

func (window *BoardWindow) HandleInput(model *Model, msg string) tea.Cmd {
	game := model.game

	switch msg {
	case "up", "k":
		if window.CursorRow == 0 {
			window.CursorRow = game.Rows() - 1
			break
		}
		window.CursorRow--

	case "down", "j":
		if window.CursorRow == game.Rows()-1 {
			window.CursorRow = 0
			break
		}
		window.CursorRow++

	case "left", "h":
		if window.CursorCol == 0 {
			window.CursorCol = game.Cols() - 1
			break
		}
		window.CursorCol--

	case "right", "l":
		if window.CursorCol == game.Cols()-1 {
			window.CursorCol = 0
			break
		}
		window.CursorCol++

	case " ", "enter":
		game.Reveal(window.CursorRow, window.CursorCol)

	case "f":
		game.ToggleFlag(window.CursorRow, window.CursorCol)

	case "r":
		model.newGame()

	case "q":
		window.CursorCol = 0
		window.CursorRow = 0

		model.newGame()
		model.CurrentWindow = model.MenuWin
	}

	return nil
}

func (window *BoardWindow) MinSize(model *Model) (int, int) {
	game := model.game

	boardWidth := game.Cols() * 2 // char + space
	minWidth := boardWidth + 4    // board + borders + padding
	minHeight := game.Rows() + 6  // title + hint + borders + padding

	return minWidth, minHeight
}

// Returns a hint string based on the current game state.
func getHint(model *Model) string {
	switch model.game.State() {
	case minesweeper.Won:
		return fmt.Sprintf("You won! Time: %s", model.game.FormattedTime())

	case minesweeper.Lost:
		return fmt.Sprintf("You lost! Time: %s", model.game.FormattedTime())

	default:
		return fmt.Sprintf("Flags: %d  Time: %s", model.game.FlagsRemaining(), model.game.FormattedTime())
	}
}
