package gentable

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
)

func New[V comparable]() Model[V] {
	m := Model[V]{
		Table: table.New().DisableOverflowRow(),
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
			if win, ok := m.data.(*window[V]); ok {
				win.PageUp()
			}
		case key.Matches(msg, keys.PageDown):
			if win, ok := m.data.(*window[V]); ok {
				win.PageDown()
			}
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
}
