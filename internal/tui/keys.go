package tui

import "charm.land/bubbles/v2/key"

// appKeys are the keys that work on every window.
var appKeys = struct {
	ForceQuit key.Binding
}{
	ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
}
