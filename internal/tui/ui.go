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
		model.width = msg.Width
		model.height = msg.Height

	case startGameMsg:
		model.board.NewGame(msg.settings)
		model, cmd := model.switchWindow(model.board)
		return model, cmd

	case exitBoardMsg:
		model, cmd := model.switchWindow(model.menu)
		return model, cmd
	}

	var cmd tea.Cmd
	model.currentWindow, cmd = model.currentWindow.Update(msg)

	return model, cmd
}

// Tea view rendering, handles terminal size checks and delegates to the current window.
func (model Model) View() tea.View {
	minWidth, minHeight := model.currentWindow.MinSize()

	var content string
	if model.width < minWidth || model.height < minHeight {
		content = fmt.Sprintf("Terminal too small — need at least %dx%d", minWidth, minHeight)
	} else {
		// Render centered base panel
		content = model.currentWindow.View()
	}

	view := tea.NewView(lg.Place(model.width, model.height, lg.Center, lg.Center, content))
	view.AltScreen = true

	return view
}
