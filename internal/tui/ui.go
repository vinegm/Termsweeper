package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

// Tea update loop
func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		return model, nil

	case time.Time:
		if model.CurrentWindow == model.BoardWin {
			return model, tick()
		}
		return model, nil

	case tea.KeyMsg:
		cmd := model.CurrentWindow.HandleInput(&model, strings.ToLower(msg.String()))
		return model, cmd
	}

	return model, nil
}

// Tea view rendering, handles terminal size checks and delegates to menu or game renderers.
func (model Model) View() string {
	minWidth, minHeight := model.CurrentWindow.MinSize(&model)

	if model.width < minWidth || model.height < minHeight {
		warn := fmt.Sprintf("Terminal too small — need at least %dx%d", minWidth, minHeight)
		return lg.Place(model.width, model.height, lg.Center, lg.Center, warn)
	}

	// Render centered base panel
	basePanel := model.CurrentWindow.Render(&model)
	return lg.Place(model.width, model.height, lg.Center, lg.Center, basePanel)
}

// tick returns a command that sends a time.Time every second to drive UI updates.
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return t })
}
