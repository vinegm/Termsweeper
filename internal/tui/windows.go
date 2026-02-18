package tui

import tea "github.com/charmbracelet/bubbletea"

// Interface that holds methods to handle rendering and input for different windows.
type Window interface {
	Render(model *Model) string
	HandleInput(model *Model, msg string) tea.Cmd
	MinSize(model *Model) (int, int)
}
