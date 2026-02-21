package tui

import tea "charm.land/bubbletea/v2"

// Window is the interface that holds methods to handle
// updates and rendering for different windows.
type Window interface {
	Update(msg tea.Msg) tea.Cmd
	View(width, height int) string
	MinSize() (int, int)
}
