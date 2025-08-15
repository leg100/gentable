package gentable

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
)

func New[V comparable](opts ...Option[V]) Model[V] {
	window := newWindow[V](10)
	m := Model[V]{
		Table:  table.New().DisableOverflowRow().Data(window),
		window: window,
	}
	for _, fn := range opts {
		fn(&m)
	}
	return m
}

type Model[V comparable] struct {
	*table.Table
	*window[V]
}

func (m Model[V]) Init() tea.Cmd {
	return nil
}

func (m Model[V]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.PageUp):
			m.window.PageUp()
		case key.Matches(msg, keys.PageDown):
			m.window.PageDown()
		}
	}
	return m, nil
}

func (m Model[V]) View() string {
	return m.Table.String()
}

func (m *Model[V]) Height(height int) {
	m.Table.Height(height)

	// TODO: we set a min of 1 because lipgloss's table has a min of 1, but we
	// should change that in the lipgloss fork.
	m.window.size = max(1, height-m.Table.NonRowHeight())

	// TODO: clamp cursor on window, maybe use a new setSize() method
}
