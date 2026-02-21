package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"

	"termsweeper/internal/tui/board"
	"termsweeper/internal/tui/menu"
)

// Tea update loop
func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		return model, nil

	case menu.StartGameMsg:
		model.board.NewGame(msg.Settings)
		model.currentWindow = model.board
		return model, nil

	case board.ExitMsg:
		model.currentWindow = model.menu
		return model, nil
	}

	return model, model.currentWindow.Update(msg)
}

// Tea view rendering, handles terminal size checks and delegates to the current window.
func (model Model) View() tea.View {
	minWidth, minHeight := model.currentWindow.MinSize()

	var content string
	if model.width < minWidth || model.height < minHeight {
		content = fmt.Sprintf("Terminal too small — need at least %dx%d", minWidth, minHeight)
	} else {
		// Render centered base panel
		content = model.currentWindow.View(model.width, model.height)
	}

	view := tea.NewView(lg.Place(model.width, model.height, lg.Center, lg.Center, content))
	view.AltScreen = true

	return view
}
