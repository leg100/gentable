package gentable

import (
	"slices"

	"github.com/charmbracelet/lipgloss/v2/table"
)

type window[V comparable] struct {
	*data[V]
	windowTable
	cursor cursor[V]
	sort   func(V, V) int
}

// cursor is the currently highlighted row
type cursor[V comparable] struct {
	n int // row position in unwindowed data
	v V
}

// windowTable provides the window with access to the underlying table.
type windowTable interface {
	YOffset(o int) *table.Table
	GetYOffset() int
	VisibleRows() int
}

func newWindow[V comparable](table windowTable) *window[V] {
	return &window[V]{
		data: &data[V]{
			cells: make(map[V][]string),
		},
		windowTable: table,
	}
}

// Append appends rows to the table data.
func (m *window[V]) Append(rows ...Row[V]) {
	m.data.add(rows...)

	if m.sort != nil {
		slices.SortFunc(m.rows(), m.sort)
	}
	m.reset()
}

func (m *window[V]) ApplyFilter(fn func(V) bool) {
	m.data.applyFilter(fn)
	m.reset()
}

func (m *window[V]) RemoveFilter() {
	m.data.removeFilter()
	m.reset()
}

func (m *window[V]) PageUp() {
	m.moveCursor(-m.VisibleRows())
}

func (m *window[V]) PageDown() {
	m.moveCursor(m.VisibleRows())
}

func (m *window[V]) CursorIndex() int {
	return m.cursor.n
}

func (m *window[V]) toTop() {
	m.moveCursor(-m.cursor.n)
}

func (m *window[V]) toBottom() {
	m.moveCursor(len(m.rows()) - m.cursor.n)
}

func (m *window[V]) moveCursor(delta int) {
	if len(m.rows()) == 0 {
		return
	}
	m.cursor.n = clamp(m.cursor.n+delta, 0, len(m.rows())-1)
	m.cursor.v = m.rows()[m.cursor.n]
	m.setYOffset()
}

// reset resets the window, re-establishing the cursor and offset; this is
// necessary whenever visible rows are re-ordered or removed.
func (m *window[V]) reset() {
	if len(m.rows()) == 0 {
		return
	}
	// Check cursor index still corresponds to cursor value
	var found bool
	if m.cursor.v != m.rows()[m.cursor.n] {
		// Value no longer corresponds, so search for value, and re-set cursor
		// index
		for i := range m.rows() {
			if m.cursor.v == m.rows()[i] {
				m.cursor.n = i
				found = true
				break
			}
		}
		if !found {
			// Value corresponding to cursor can not be found; this happens when the
			// cursor has not been set yet or the value has been filtered out.
			m.cursor.v = m.rows()[0]
			m.cursor.n = 0
		}
	}
	m.setYOffset()
}

func (m *window[V]) setYOffset() {
	// Offset must be at least the cursor index minus the max number
	// of visible rows.
	minimum := max(0, m.cursor.n-m.VisibleRows()+1)
	// Offset must be at most the lesser of:
	// (a) the cursor index, or
	// (b) the number of rows minus the maximum number of visible rows (so that
	// as many rows as possible are rendered)
	maximum := max(0, min(m.cursor.n, len(m.rows())-m.VisibleRows()))
	_ = m.YOffset(clamp(m.GetYOffset(), minimum, maximum))
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
