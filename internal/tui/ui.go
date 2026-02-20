package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"
)

// Tea update loop
func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.screen.Width = msg.Width
		model.screen.Height = msg.Height
		return model, nil

	case startGameMsg:
		model.board.NewGame(msg.settings)
		model.currentWindow = model.board
		return model, nil

	case exitBoardMsg:
		model.currentWindow = model.menu
		return model, nil
	}

	var cmd tea.Cmd
	model.currentWindow, cmd = model.currentWindow.Update(msg)

	return model, cmd
}

// Tea view rendering, handles terminal size checks and delegates to the current window.
func (model Model) View() tea.View {
	minWidth, minHeight := model.currentWindow.MinSize()

	var content string
	if model.screen.Width < minWidth || model.screen.Height < minHeight {
		content = fmt.Sprintf("Terminal too small — need at least %dx%d", minWidth, minHeight)
	} else {
		// Render centered base panel
		content = model.currentWindow.View()
	}

	view := tea.NewView(lg.Place(model.screen.Width, model.screen.Height, lg.Center, lg.Center, content))
	view.AltScreen = true

	return view
}
