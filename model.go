package gentable

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
)

func NewDefault() Model[[]string] {
	var id int
	getID := func(v []string) ID {
		newID := id
		id++
		return newID
	}
	render := func(v []string) []string { return v }
	m := Model[[]string]{
		data: NewData[[]string](getID, render),
		lt:   table.New().DisableOverflowRow(),
	}
	return m
}

type Model[V any] struct {
	lt   *table.Table
	data *data[V]
}

func (m Model[V]) Init() tea.Cmd {
	return nil
}

func (m Model[V]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//switch msg := msg.(type) {
	//case tea.KeyMsg:
	//	switch {
	//	case key.Matches(msg, keys.PageUp):
	//		m.data.PageUp()
	//	}
	//}
	return m, nil
}

func (m Model[V]) View() string {
	return m.lt.String()
}

func (m *Model[V]) Height(height int) {
	m.lt.Height(height)
	if window, ok := m.data.(interface {
		SetWindowSize(int)
	}); ok {
		window.SetWindowSize(max(1, height-m.lt.NonRowHeight()))
	}
}
