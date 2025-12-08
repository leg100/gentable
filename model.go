package gentable

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

const headerHeight = 1

func New[V comparable](headers ...string) Model[V] {
	m := Model[V]{
		focused: true,
		headers: headers,
	}
	m.window = newWindow[V](headers...)
	return m
}

type Model[V comparable] struct {
	*window[V]
	headers []string
	focused bool
	debug   bool
}

type StyleFunc func(row, col int) lipgloss.Style

func (m Model[V]) Init() tea.Cmd {
	return nil
}

func (m Model[V]) Update(msg tea.Msg) (Model[V], tea.Cmd) {
	if !m.focused {
		return m, nil
	}
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
		case key.Matches(msg, keys.Select):
			m.window.data.toggleSelect(m.cursor.n)
		case key.Matches(msg, keys.SelectAll):
			m.window.data.selectAll()
		case key.Matches(msg, keys.SelectClear):
			m.window.data.clearSelection()
		}
	}
	return m, nil
}

func (m Model[V]) View() string {
	var sb strings.Builder
	if m.debug {
		sb.WriteString(m.constructDebugHeader())
	}
	sb.WriteString(m.constructHeaders())

	// If there are no data rows render nothing.
	if m.data.Rows() > 0 {
		for r := m.window.offset; r < m.window.offset+m.visibleRows(); r++ {
			sb.WriteRune('\n')
			sb.WriteString(m.constructRow(r))
		}
	}

	return sb.String()
}

func (m Model[V]) constructDebugHeader() string {
	return fmt.Sprintf("table_width: %d; window_size: %d; visible: %d; offset: %d; cursor: %d\n",
		m.resizer.tableWidth,
		m.window.size,
		m.window.visibleRows(),
		m.window.offset,
		m.window.cursor.n,
	)
}

// Height sets the table height.
func (m *Model[V]) Height(height int) {
	if m.debug {
		height -= 1
	}
	m.window.size = height - headerHeight
	// Altering the window size necessitates resetting the window, to check if
	// offset or cursor needs moving.
	m.window.reset()
}

func (m *Model[V]) Focus(focus bool) {
	m.focused = focus
}

// HeaderRow denotes the header's row index used when rendering headers. Use
// this value when looking to customize header styles in StyleFunc.
const HeaderRow int = -1

// constructHeaders constructs the headers for the table given it's current
// header configuration and data.
func (t *Model[V]) constructHeaders() string {
	var s strings.Builder
	// TODO: why this calculation?
	cells := make([]string, 0, len(t.headers)*2+1)

	for j, header := range t.headers {
		header = t.truncateCell(header, HeaderRow, j)

		cells = append(cells,
			lipgloss.NewStyle().
				Height(1).
				Width(t.columnWidths[j]).
				Render(header),
		)
	}

	for i, cell := range cells {
		cells[i] = strings.TrimRight(cell, "\n")
	}

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cells...))

	return s.String()
}

// constructRow constructs the row for the table given an index and row data
// based on the current configuration. If isOverflow is true, the row is
// rendered as an overflow row (using ellipsis).
func (t *Model[V]) constructRow(index int) string {
	var s strings.Builder
	// TODO: why this calculation?
	cells := make([]string, 0, t.data.Columns()*2+1)

	for c := range t.data.Columns() {
		cell := t.data.At(index, c)
		cellStyle := t.style(index, c)
		cell = t.truncateCell(cell, index, c)
		cells = append(cells, cellStyle.
			// Account for the margins in the cell sizing.
			Height(1).
			MaxHeight(1).
			Width(t.columnWidths[c]-cellStyle.GetHorizontalMargins()).
			MaxWidth(t.columnWidths[c]).
			Render(cell))
	}

	for i, cell := range cells {
		cells[i] = strings.TrimRight(cell, "\n")
	}

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cells...))

	return s.String()
}

func (t *Model[V]) truncateCell(cell string, rowIndex, colIndex int) string {
	cellWidth := t.columnWidths[colIndex]
	cellStyle := t.style(rowIndex, colIndex)

	length := cellWidth - cellStyle.GetHorizontalPadding() - cellStyle.GetHorizontalMargins()
	return ansi.Truncate(cell, length, "…")
}

func (t *Model[V]) SetSortFunc(fn func(V, V) int) {
	t.sort = fn
}

func (t *Model[V]) EnableDebug() {
	t.debug = true
}

// style returns the style for a cell based on its position (row, column).
func (m *Model[V]) style(row, _ int) lipgloss.Style {
	if m.disableCursorHighlighting {
		return lipgloss.NewStyle()
	}
	s := lipgloss.NewStyle()
	switch {
	case row == m.cursor.n:
		return s.
			Foreground(defaultCursorForeground).
			Background(defaultCursorBackground)
	case m.data.isSelected(row):
		return s.
			Background(defaultSelectionBackground)
	default:
		return s
	}
}
