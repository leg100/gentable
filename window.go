package gentable

import (
	"slices"
)

// window is a "window" of the underlying data, showing only the visible rows
// including the cursor. If the user moves the cursor then the window "tracks"
// the cursor, changing the offset accordingly.
//
// The cursor tracks the row it is set to. If sorting is enabled and a row is
// appended, then the window offset may change to ensure the cursor is still
// visible. If many rows are being appended this can result in the cursor and
// thus the window "jumping around". To mitigate this, at startup, if the user
// has not moved the cursor, the cursor will instead track the first row.
type window[V comparable] struct {
	*data[V]
	cursor  cursor[V]
	sort    func(V, V) int
	offset  int
	size    int
	startup bool

	// for testing purposes
	disableCursorHighlighting bool
}

// cursor is the currently highlighted row
type cursor[V comparable] struct {
	n int // row position in unwindowed data
	v V
}

func newWindow[V comparable](headers ...string) *window[V] {
	w := &window[V]{
		startup: true,
		cursor:  cursor[V]{},
	}
	w.data = newData[V](headers...)
	return w
}

// Append appends rows to the table data.
func (m *window[V]) Append(rows ...Row[V]) {
	m.data.add(rows...)

	if m.sort != nil {
		slices.SortFunc(m.rows(), m.sort)
	}
	m.reset()
}

// DeleteAll deletes all table rows.
func (w *window[V]) DeleteAll() {
	w.data = newData[V](w.data.headers...)
	w.cursor = cursor[V]{}
	w.reset()
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
	m.moveCursor(-m.visibleRows())
}

func (m *window[V]) PageDown() {
	m.moveCursor(m.visibleRows())
}

func (m *window[V]) CursorIndex() int {
	return m.cursor.n
}

func (m *window[V]) CursorValue() V {
	return m.cursor.v
}

func (m *window[V]) toTop() {
	m.moveCursor(-m.cursor.n)
}

func (m *window[V]) toBottom() {
	m.moveCursor(len(m.rows()) - m.cursor.n)
}

func (m *window[V]) visibleRows() int {
	return min(m.data.Rows()-m.offset, m.size)
}

func (m *window[V]) moveCursor(delta int) {
	if len(m.rows()) == 0 {
		return
	}
	m.startup = false
	m.cursor.n = clamp(m.cursor.n+delta, 0, len(m.rows())-1)
	m.cursor.v = m.rows()[m.cursor.n]
	m.setOffset()
}

// reset resets the window, re-establishing the cursor and offset; this is
// necessary whenever visible rows are re-ordered or removed.
func (m *window[V]) reset() {
	if len(m.rows()) == 0 {
		return
	}
	if m.startup {
		m.cursor.v = m.rows()[0]
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
	m.setOffset()
}

func (m *window[V]) setOffset() {
	// Offset must at the very least be set to a value that allows the cursor to
	// be on the last visible row.
	minimum := max(0, m.cursor.n-m.visibleRows()+1)
	// Offset must be at most the lesser of:
	// (a) the cursor index, or
	// (b) the number of rows minus the maximum number of visible rows (so that
	// as many rows as possible are rendered)
	maximum := max(0, min(m.cursor.n, len(m.rows())-m.visibleRows()))
	m.offset = clamp(m.offset, minimum, maximum)
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
