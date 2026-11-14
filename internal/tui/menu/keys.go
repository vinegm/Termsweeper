package menu

import "charm.land/bubbles/v2/key"

// keyMap is the menu's keys.
type keyMap struct {
	Up, Down, Left, Right, Select, Quit key.Binding

	// SizeDone and Cancel are the keys while resizing the
	// board, SeedDone and Delete while typing the seed.
	SizeDone, Cancel, SeedDone, Delete key.Binding

	// Move, Change and Resize name the keys
	// above as one in the key hints bar.
	Move, Change, Resize key.Binding
}

var keys = keyMap{
	Up:    key.NewBinding(key.WithKeys("up", "k")),
	Down:  key.NewBinding(key.WithKeys("down", "j")),
	Left:  key.NewBinding(key.WithKeys("left", "h")),
	Right: key.NewBinding(key.WithKeys("right", "l")),
	Select: key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter", "Select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "Quit"),
	),
	SizeDone: key.NewBinding(
		key.WithKeys("enter", "space"),
		key.WithHelp("enter", "Done"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Cancel"),
	),
	SeedDone: key.NewBinding(
		key.WithKeys("enter", "esc"),
		key.WithHelp("enter", "Done"),
	),
	Delete: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "Delete"),
	),
	Move: key.NewBinding(
		key.WithKeys("down", "j", "up", "k"),
		key.WithHelp("jk", "Move"),
	),
	Change: key.NewBinding(
		key.WithKeys("left", "h", "right", "l"),
		key.WithHelp("hl", "Change"),
	),
	Resize: key.NewBinding(
		key.WithKeys("left", "h", "down", "j", "up", "k", "right", "l"),
		key.WithHelp("hjkl", "Resize"),
	),
}

// lowered is a key by its lowercased name,
// so the keys work the same with caps lock on.
type lowered string

// String is the name, for key.Matches.
func (name lowered) String() string {
	return string(name)
}
