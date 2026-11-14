package board

import "charm.land/bubbles/v2/key"

// keyMap is the board's keys.
type keyMap struct {
	Up, Down, Left, Right key.Binding

	Reveal, Flag, NewGame, Menu key.Binding

	// Move names the keys above as
	// one in the key hints bar.
	Move key.Binding
}

var keys = keyMap{
	Up:    key.NewBinding(key.WithKeys("up", "k", "w")),
	Down:  key.NewBinding(key.WithKeys("down", "j", "s")),
	Left:  key.NewBinding(key.WithKeys("left", "h", "a")),
	Right: key.NewBinding(key.WithKeys("right", "l", "d")),
	Reveal: key.NewBinding(
		key.WithKeys("space", "enter"),
		key.WithHelp("space", "Reveal"),
	),
	Flag: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "Flag"),
	),
	NewGame: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "New game"),
	),
	Menu: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "Menu"),
	),
	Move: key.NewBinding(
		key.WithKeys("left", "h", "a", "down", "j", "s", "up", "k", "w", "right", "l", "d"),
		key.WithHelp("hjkl", "Move"),
	),
}

// lowered is a key by its lowercased name,
// so the keys work the same with caps lock on.
type lowered string

// String is the name, for key.Matches.
func (name lowered) String() string {
	return string(name)
}
