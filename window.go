package gentable

import "slices"

type window[V comparable] struct {
	*data[V]

	cursor cursor[V]

	// start is the unwindowed row index of the first visible row.
	start int
	// size is the number of terminal rows the window occupies (not necessarily
	// populated with data). Zero means it occupies an unlimited number of
	// terminal rows, effectively disabling windowing.
	size int

	sort func(V, V) int
}

// cursor is the currently highlighted row
type cursor[V comparable] struct {
	n int // row position in unwindowed data
	v V
}

func newWindow[V comparable](size int) *window[V] {
	return &window[V]{
		size: size,
		data: &data[V]{
			cells: make(map[V][]string),
		},
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

// At returns the contents of the cell at the given index of windowed data.
func (m *window[V]) At(row, cell int) string {
	v := m.rows()[row+m.start]
	cells := m.cells[v]
	if cell >= len(cells) {
		// Not all rows have the same number of cells but the lipgloss table lib
		// that calls this method doesn't know that.
		return ""
	}
	return cells[cell]
}

// Rows returns the number of populated, visible rows.
func (m *window[V]) Rows() int {
	if m.size == 0 {
		return len(m.rows())
	}
	return min(m.size, len(m.rows()))
}

func (m *window[V]) applyFilter(fn func(V) bool) {
	m.data.applyFilter(fn)
	m.reset()
}

func (m *window[V]) removeFilter() {
	m.data.removeFilter()
	m.reset()
}

func (m *window[V]) PageUp() {
	m.moveCursor(-m.size)
}

func (m *window[V]) PageDown() {
	m.moveCursor(m.size)
}

func (m *window[V]) StartIndex() int {
	return m.start
}

func (m *window[V]) CursorIndex() int {
	return m.cursor.n
}

func (m *window[V]) Size() int {
	return m.size
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
	m.setStart()
}

// reset resets the window, re-establishing the cursor and start rows; this is
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
	}
	if !found {
		// Value corresponding to cursor can not be found; this happens when the
		// cursor has not been set yet or the value has been filtered out.
		m.cursor.v = m.rows()[0]
		m.cursor.n = 0
	}
	m.setStart()
}

func (m *window[V]) setStart() {
	// Start index must be at least the cursor index minus the max number
	// of visible rows.
	minimum := max(0, m.cursor.n-m.size+1)
	// Start index must be at most the lesser of:
	// (a) the cursor index, or
	// (b) the number of rows minus the maximum number of visible rows (so that
	// as many rows as possible are rendered)
	maximum := max(0, min(m.cursor.n, len(m.rows())-m.size))
	m.start = clamp(m.start, minimum, maximum)
}

func (m *window[V]) setSize(size int) {
	m.size = size
	m.reset()
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
