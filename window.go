package gentable

import (
	"slices"
)

// window is the viewport of visible rows
type window[V comparable] struct {
	unwindowed data[V]
	cursor     cursor[V]
	// start is the unwindowed row index of the first visible row.
	start int
	// size is the number of terminal rows the window occupies (not necessarily
	// populated with data). Zero means it occupies an unlimited number of
	// terminal rows, effectively disabling windowing.
	size int
	sort func(V, V) int
}

func newWindow[V comparable](size int) *window[V] {
	return &window[V]{
		size: size,
		unwindowed: &base[V]{
			cells: make(map[V][]string),
		},
	}
}

// cursor is the currently highlighted row
type cursor[V comparable] struct {
	n int // row position in unwindowed data
	v V
}

// At returns the contents of the cell at the given index of windowed data.
func (m *window[V]) At(row, cell int) string {
	// TODO: this should never happen?
	//if row >= m.size {
	//	return ""
	//}
	//// TODO: this should never happen?
	//if m.start+row >= len(m.unwindowed.Rows()) {
	//	return ""
	//}
	v := m.unwindowed.Rows()[row+m.start]
	cells := m.unwindowed.getCells(v)
	if cell >= len(cells) {
		// Not all rows have the same number of cells but the lipgloss table lib
		// that calls this method doesn't know that.
		return ""
	}
	return cells[cell]
}

// Rows returns the number of populated rows in the window
func (m *window[V]) Rows() int {
	if m.size == 0 {
		return len(m.unwindowed.Rows())
	}
	return min(m.size, len(m.unwindowed.Rows()))
}

func (m *window[V]) Columns() int {
	return m.unwindowed.Columns()
}

func (m *window[V]) Append(rows ...row[V]) {
	m.unwindowed.Append(rows...)

	if m.sort != nil {
		slices.SortFunc(m.unwindowed.Rows(), m.sort)
	}
	m.reset()
}

// filter applies a filter to the unwindowed rows. If fn is nil and a filter is
// currently applied then it is removed.
func (m *window[V]) filter(fn func(V) bool) {
	if filter, ok := m.unwindowed.(*filter[V]); ok {
		if fn == nil {
			// Remove filter
			m.unwindowed = filter.data
			return
		}
	}
	m.unwindowed = newFilter(m.unwindowed, fn)
	m.reset()
}

func (m *window[V]) PageUp() {
	m.moveCursor(-m.size)
}

func (m *window[V]) PageDown() {
	m.moveCursor(+m.size)
}

func (m *window[V]) moveCursor(delta int) {
	if len(m.unwindowed.Rows()) == 0 {
		return
	}
	m.cursor.n = clamp(m.cursor.n+delta, 0, len(m.unwindowed.Rows())-1)
	m.cursor.v = m.unwindowed.Rows()[m.cursor.n]
	m.setStart()
}

// reset resets the window, re-establishing the cursor and start rows; this is
// necessary whenever rows are re-ordered or removed.
func (m *window[V]) reset() {
	// Check cursor index still corresponds to cursor value
	var found bool
	if m.cursor.v != m.unwindowed.Rows()[m.cursor.n] {
		// Value no longer corresponds, so search for value, and re-set cursor
		// index
		for i := range m.unwindowed.Rows() {
			if m.cursor.v == m.unwindowed.Rows()[i] {
				m.cursor.n = i
				found = true
				break
			}
		}
	}
	if !found {
		// Value corresponding to cursor can not be found; this happens when the
		// cursor has not been set yet or the value has been filtered out.
		if len(m.unwindowed.Rows()) > 0 {
			m.cursor.v = m.unwindowed.Rows()[0]
		}
	}
	m.setStart()
}

func (m *window[V]) setStart() {
	// Start index must be at least the cursor index minus the max number
	// of visible rows.
	minimum := max(0, m.cursor.n-m.size)
	// Start index must be at most the lesser of:
	// (a) the cursor index, or
	// (b) the number of rows minus the maximum number of visible rows (so that
	// as many rows as possible are rendered)
	maximum := max(0, min(m.cursor.n, len(m.unwindowed.Rows())-m.size))
	m.start = clamp(m.start, minimum, maximum)
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
