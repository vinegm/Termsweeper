package theme

import (
	"fmt"

	lg "charm.land/lipgloss/v2"

	"github.com/vinegm/termsweeper/internal/minesweeper"
)

// CellStyle is the style of a cell, by its state and whether it is selected or highlighted.
func (theme *Theme) CellStyle(cell minesweeper.Cell, isSelected, isHighlighted bool) lg.Style {
	if isSelected {
		return theme.selectedSquareStyle
	}

	var style lg.Style

	switch {
	case cell.StateIs(minesweeper.Exploded):
		style = theme.explodedSquareStyle

	case cell.StateIs(minesweeper.Revealed) && cell.StateIs(minesweeper.Mined):
		style = theme.minedSquareStyle

	case cell.StateIs(minesweeper.Flagged):
		style = theme.flaggedSquareStyle

	case cell.StateIs(minesweeper.Revealed):
		style = theme.hintStyle(cell)

	case cell.StateIs(minesweeper.Odd):
		style = theme.oddSquareStyle

	default:
		style = theme.evenSquareStyle
	}

	style = theme.handleGridBg(style, cell)

	if isHighlighted {
		style = withBg(style, theme.highlightColor(cell))
	}

	return style
}

// CellChar is the character a cell shows in its state.
func (theme *Theme) CellChar(cell minesweeper.Cell) string {
	if cell.StateIs(minesweeper.Revealed) {
		if cell.StateIs(minesweeper.Mined) {
			return theme.config.Chars.Mine
		}

		if cell.Adj == 0 {
			return " "
		}

		return fmt.Sprintf("%d", cell.Adj)
	}

	if cell.StateIs(minesweeper.Flagged) {
		return theme.config.Chars.Flag
	}

	return theme.config.Chars.Unrevealed
}

// handleGridBg Restores the checkerboard background on a square whose
// kind is configured to keep it, and leaves every other square alone.
func (theme *Theme) handleGridBg(style lg.Style, cell minesweeper.Cell) lg.Style {
	preserve := theme.config.Theme.PreserveGridBg

	var keep bool

	switch {
	case cell.StateIs(minesweeper.Exploded) || (cell.StateIs(minesweeper.Revealed) && cell.StateIs(minesweeper.Mined)):
		keep = preserve.Mine

	case cell.StateIs(minesweeper.Flagged):
		keep = preserve.Flag

	case cell.StateIs(minesweeper.Revealed):
		keep = preserve.Revealed
	}

	if !keep {
		return style
	}

	if cell.StateIs(minesweeper.Odd) {
		return style.Background(theme.oddSquareStyle.GetBackground())
	}

	return style.Background(theme.evenSquareStyle.GetBackground())
}

// Returns the style for a revealed square, by its adjacent mine count.
func (theme *Theme) hintStyle(cell minesweeper.Cell) lg.Style {
	if int(cell.Adj) < len(theme.config.Theme.Squares.Hints) {
		return theme.squareMineHintColor[cell.Adj]
	}

	return theme.TextStyle
}

// Returns the background a highlighted square takes, keeping the checkerboard.
func (theme *Theme) highlightColor(cell minesweeper.Cell) string {
	if cell.StateIs(minesweeper.Odd) {
		return theme.config.Theme.Highlight.Odd
	}

	return theme.config.Theme.Highlight.Even
}
