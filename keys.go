package gentable

import (
	"github.com/charmbracelet/bubbles/key"
)

var keys = struct {
	LineUp     key.Binding
	LineDown   key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	GotoTop    key.Binding
	GotoBottom key.Binding
}{
	LineUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	LineDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
	GotoTop: key.NewBinding(
		key.WithKeys("home", "g"),
		key.WithHelp("g/home", "go to start"),
	),
	GotoBottom: key.NewBinding(
		key.WithKeys("end", "G"),
		key.WithHelp("G/end", "go to end"),
	),
}
