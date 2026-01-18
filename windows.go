package main

import tea "github.com/charmbracelet/bubbletea"

// Interface that holds methods to handle rendering and input for different windows.
type Window interface {
	Render(model *model) string
	HandleInput(model *model, msg tea.Msg) tea.Cmd
	MinSize() (int, int)
}
