package gentable

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func New[V comparable](opts ...Option[V]) Model[V] {
	window := newWindow[V](10)
	m := Model[V]{
		rendering: table.New().DisableOverflowRow().Data(window),
		window:    window,
	}
	for _, fn := range opts {
		fn(&m)
	}
	return m
}

type Model[V comparable] struct {
	rendering *table.Table
	*window[V]
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
		}
	}
	m.rendering.StyleFunc(func(row, col int) lipgloss.Style {
		if m.window.start+row == m.cursor.n {
			return lipgloss.NewStyle().Background(lipgloss.Color("#ffffff"))
		}
		return lipgloss.NewStyle()
	})
	return m, nil
}

func (m Model[V]) View() string {
	return m.rendering.String()
}

func (m *Model[V]) Filter(fn func(V) bool) {
	m.window.applyFilter(fn)
}

func (m *Model[V]) RemoveFilter() {
	m.window.removeFilter()
}

// Height sets the table height.
func (m *Model[V]) Height(height int) {
	m.rendering.Height(height)

	// TODO: we set a minimum of 1 because lipgloss's table has a minimum of 1,
	// but we should change that in the lipgloss fork.
	m.window.setSize(max(1, height-m.rendering.NonRowHeight()))
}

// Headers sets the table headers.
func (t *Model[V]) Headers(headers ...string) *Model[V] {
	t.rendering.Headers(headers...)
	return t
}
