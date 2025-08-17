package gentable

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/table"
)

func New[V comparable](opts ...Option[V]) Model[V] {
	table := table.New().Overflow(false)
	window := newWindow[V](table)
	table.Data(window)

	m := Model[V]{
		table:  table,
		window: window,
	}
	for _, fn := range opts {
		fn(&m)
	}
	return m
}

type Model[V comparable] struct {
	*window[V]
	// table is the wrapped lipgloss table responsible for actual rendering
	table *table.Table
}

func (m Model[V]) Init() tea.Cmd {
	return nil
}

func (m Model[V]) Update(msg tea.Msg) (Model[V], tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.LineUp):
			m.window.moveCursor(-1)
		case key.Matches(msg, keys.LineDown):
			m.window.moveCursor(1)
		case key.Matches(msg, keys.PageUp):
			m.window.PageUp()
		case key.Matches(msg, keys.PageDown):
			m.window.PageDown()
		case key.Matches(msg, keys.GotoTop):
			m.window.toTop()
		case key.Matches(msg, keys.GotoBottom):
			m.window.toBottom()
		case key.Matches(msg, keys.GotoTop):
			m.window.toTop()
		case key.Matches(msg, keys.GotoBottom):
			m.window.toBottom()
		case key.Matches(msg, keys.Select):
			m.window.data.toggleSelect(m.cursor.n)
		case key.Matches(msg, keys.SelectAll):
			m.window.data.selectAll()
		case key.Matches(msg, keys.SelectClear):
			m.window.data.clearSelection()
		}
	}
	m.table.StyleFunc(func(row, col int) lipgloss.Style {
		s := lipgloss.NewStyle()
		if row == m.cursor.n {
			return s.Background(lipgloss.Color("#ffffff"))
		}
		if m.window.data.isSelected(row) {
			return s.Background(lipgloss.Color("#DBBD70"))
		}
		return s
	})
	return m, nil
}

func (m Model[V]) View() string {
	return m.table.String()
}

// Height sets the table height.
func (m *Model[V]) Height(height int) {
	m.table.Height(height)
	// Setting the height alters the number of visible rows, so the window needs
	// resetting.
	m.window.reset()
}

// Width sets the table width.
func (m *Model[V]) Width(width int) {
	m.table.Width(width)
}

// Headers sets the table headers.
func (m *Model[V]) Headers(headers ...string) *Model[V] {
	m.table.Headers(headers...)
	// adding headers can alter the number of visible rows, so the window needs
	// resetting.
	m.window.reset()
	return m
}

func (m *Model[V]) FirstVisibleRowIndex() int {
	return m.table.FirstVisibleRowIndex()
}

func (m *Model[V]) LastVisibleRowIndex() int {
	return m.table.LastVisibleRowIndex()
}

// Border sets the table border.
func (m *Model[V]) Border(border lipgloss.Border) *Model[V] {
	m.table.Border(border)
	return m
}

// Wrap dictates whether or not the table content should wrap.
//
// This only applies to data cells. Headers are never wrapped.
func (m *Model[V]) Wrap(w bool) *Model[V] {
	m.table.Wrap(w)
	return m
}
