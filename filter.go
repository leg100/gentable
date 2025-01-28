package gentable

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Height of filter widget
const filterHeight = 2

type filter struct {
	textinput.Model
}

func newFilter() *filter {
	m := textinput.New()
	m.Prompt = "Filter: "
	return &filter{
		Model: m,
	}
}

func (m *filter) Update(msg tea.Msg) (bool, tea.Cmd) {
	if m.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.UnfocusFilter):
				m.Blur()
				return false, nil
			case key.Matches(msg, keys.CloseFilter):
				m.Blur()
				m.SetValue("")
				return true, nil
			default:
				var cmd tea.Cmd
				m.Model.Update(msg)
				return true, cmd
			}
		}
		var cmd tea.Cmd
		m.Model, cmd = m.Model.Update(msg)
		return false, cmd
	}
	return false, nil
}

func (m *filter) isVisible() bool {
	// Filter is visible if it's either in focus, or it has a non-empty value.
	return m.Focused() || m.Value() != ""
}
