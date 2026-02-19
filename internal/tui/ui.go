package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"
)

// Tea update loop
func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		return model, nil

	case time.Time:
		return model, tick()

	case tea.KeyPressMsg:
		cmd := model.CurrentWindow.HandleInput(&model, strings.ToLower(msg.String()))
		return model, cmd
	}

	return model, nil
}

// Tea view rendering, handles terminal size checks and delegates to menu or game renderers.
func (model Model) View() tea.View {
	minWidth, minHeight := model.CurrentWindow.MinSize(&model)

	var content string
	if model.width < minWidth || model.height < minHeight {
		content = fmt.Sprintf("Terminal too small — need at least %dx%d", minWidth, minHeight)
	} else {
		// Render centered base panel
		content = model.CurrentWindow.Render(&model)
	}

	view := tea.NewView(lg.Place(model.width, model.height, lg.Center, lg.Center, content))
	view.AltScreen = true

	return view
}

// tick returns a command that sends a time.Time every second to drive UI updates.
// The chain has to keep running even off the board screen: the renderer only
// repaints in response to messages, so dropping it would freeze the game clock.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return t })
}
