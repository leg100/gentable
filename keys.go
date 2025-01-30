package gentable

import (
	"github.com/charmbracelet/bubbles/key"
)

var keys = struct {
	PageUp   key.Binding
	PageDown key.Binding
}{
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),
}
