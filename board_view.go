package main

import (
	"fmt"
	"strings"
)

// Renders the game board.
func renderGame(model model) string {
	var sb strings.Builder

	sb.WriteString("Termsweeper\n")
	sb.WriteString(getHint(model) + "\n\n")

	for row := range model.game.board {
		for col := range model.game.board[row] {
			cell := model.game.board[row][col]
			char := cell.char()

			isCursorRow := model.game.cursorRow == row
			isCursorCol := model.game.cursorCol == col
			isSelectedCell := isCursorRow && isCursorCol
			style := cellStyle(row, col, isSelectedCell)

			sb.WriteString(style.Render(char + " "))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// Provides a hint based on the current game state.
func getHint(model model) string {
	switch model.game.state {
	case won:
		return "You won! Press q to quit."
	case lost:
		return "You lost! Press q to quit."

	default:
		return fmt.Sprintf("Flags: %d", model.game.flagsRemaining())
	}
}

