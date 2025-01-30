package gentable

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
)

func New() Model {
	data := newData()
	m := Model{
		data: data,
		lt: table.New().
			Data(data).
			DisableOverflowRow(),
	}
	return m
}

type Model struct {
	lt   *table.Table
	data *data
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.PageUp):
			m.data.PageUp()
		}
	}
	return m, nil
}

func (m Model) View() string {
	return m.lt.String()
}

func (m *Model) Rows(rows ...[]string) {
	for _, row := range rows {
		m.data.Append(row)
	}
}

func (m *Model) Height(height int) {
	m.lt.Height(height)
	m.data.size = max(1, height-m.lt.NonRowHeight())
}
