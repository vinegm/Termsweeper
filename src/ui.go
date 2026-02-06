package src

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

// Tea update loop
func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		return model, nil

	case tea.KeyMsg:
		cmd := model.CurrentWindow.HandleInput(&model, msg)
		return model, cmd
	}

	return model, nil
}

// Tea view rendering, handles terminal size checks and delegates to menu or game renderers.
func (model model) View() string {
	var minWidth, minHeight int
	minWidth, minHeight = model.CurrentWindow.MinSize(&model)

	if model.width > 0 && model.height > 0 && (model.width < minWidth || model.height < minHeight) {
		warn := fmt.Sprintf("Terminal too small — need at least %dx%d", minWidth, minHeight)
		return lg.Place(model.width, model.height, lg.Center, lg.Center, warn)
	}

	var basePanel string
	if model.CurrentWindow != nil {
		basePanel = model.CurrentWindow.Render(&model)
	}

	// Render centered base panel
	if model.width > 0 && model.height > 0 {
		centeredBase := lg.Place(model.width, model.height, lg.Center, lg.Center, basePanel)
		return centeredBase
	}

	return basePanel
}
