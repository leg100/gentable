package gentable

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
)

func New(data table.Data) Model {
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
	data table.Data
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//switch msg := msg.(type) {
	//case tea.KeyMsg:
	//	switch {
	//	case key.Matches(msg, keys.PageUp):
	//		m.data.PageUp()
	//	}
	//}
	return m, nil
}

func (m Model) View() string {
	return m.lt.String()
}

func (m *Model) Height(height int) {
	m.lt.Height(height)
	if window, ok := m.data.(interface {
		SetWindowSize(int)
	}); ok {
		window.SetWindowSize(max(1, height-m.lt.NonRowHeight()))
	}
}
