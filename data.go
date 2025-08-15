package gentable

import (
	"slices"
)

type base[V comparable] struct {
	unfiltered []V
	filtered   []V
	filter     func(V) bool
}

func (b *base[V]) getValue(i int) V {
	if b.filter != nil {
		return b.filtered[i]
	}
	return b.unfiltered[i]
}

// numRows returns the number of current rows
func (b *base[V]) numRows() int {
	if b.filter != nil {
		return len(b.filtered)
	}
	return len(b.unfiltered)
}

type data[V comparable] struct {
	rows    []V
	cells   map[V][]string
	columns int

	cursor cursor[V]
	// start is the unwindowed row index of the first visible row.
	start int
	// size is the number of terminal rows the window occupies (not necessarily
	// populated with data). Zero means it occupies an unlimited number of
	// terminal rows, effectively disabling windowing.
	size int

	sort func(V, V) int

	filtered []V
	filter   func(V) bool
}

// cursor is the currently highlighted row
type cursor[V comparable] struct {
	n int // row position in unwindowed data
	v V
}

type row[V comparable] struct {
	v     V
	cells []string
}

func newData[V comparable](size int) *data[V] {
	return &data[V]{
		size:  size,
		cells: make(map[V][]string),
	}
}

func (m *data[V]) Append(rows ...row[V]) {
	for _, row := range rows {
		m.rows = append(m.rows, row.v)
		m.cells[row.v] = row.cells
		m.columns = max(m.columns, len(row.cells))

		if m.filter != nil {
			if m.filter(row.v) {
				m.filtered = append(m.filtered, row.v)
			}
		}
	}

	if m.sort != nil {
		slices.SortFunc(m.rows, m.sort)
	}
	m.reset()
}

// At returns the contents of the cell at the given index of windowed data.
func (m *data[V]) At(row, cell int) string {
	v := m.getValue(row + m.start)
	cells := m.cells[v]
	if cell >= len(cells) {
		// Not all rows have the same number of cells but the lipgloss table lib
		// that calls this method doesn't know that.
		return ""
	}
	return cells[cell]
}

// Rows returns the number of populated, visible rows.
func (m *data[V]) Rows() int {
	var rows int
	if m.filter != nil {
		rows = len(m.filtered)
	} else {
		rows = len(m.rows)
	}
	if m.size == 0 {
		return rows
	}
	return min(m.size, rows)
}

func (m *data[V]) applyFilter(fn func(V) bool) {
	for _, v := range m.rows {
		if fn(v) {
			m.filtered = append(m.filtered, v)
		}
	}
	m.reset()
}

func (m *data[V]) removeFilter() {
	m.filter = nil
	m.reset()
}

func (m *data[V]) PageUp() {
	m.moveCursor(-m.size)
}

func (m *data[V]) PageDown() {
	m.moveCursor(+m.size)
}

func (m *data[V]) moveCursor(delta int) {
	if m.Rows() == 0 {
		return
	}
	m.cursor.n = clamp(m.cursor.n+delta, 0, m.Rows()-1)
	m.cursor.v = m.getValue(m.cursor.n)
	m.setStart()
}

func (m *data[V]) getValue(i int) V {
	if m.filter != nil {
		return m.filtered[i]
	}
	return m.rows[i]
}

// reset resets the window, re-establishing the cursor and start rows; this is
// necessary whenever rows are re-ordered or removed.
func (m *data[V]) reset() {
	// Check cursor index still corresponds to cursor value
	var found bool
	if m.cursor.v != m.rows[m.cursor.n] {
		// Value no longer corresponds, so search for value, and re-set cursor
		// index
		for i := range m.Rows() {
			if m.cursor.v == m.getValue(i) {
				m.cursor.n = i
				found = true
				break
			}
		}
	}
	if !found {
		// Value corresponding to cursor can not be found; this happens when the
		// cursor has not been set yet or the value has been filtered out.
		if m.Rows() > 0 {
			m.cursor.v = m.getValue(0)
		}
	}
	m.setStart()
}

func (m *data[V]) setStart() {
	// Start index must be at least the cursor index minus the max number
	// of visible rows.
	minimum := max(0, m.cursor.n-m.size)
	// Start index must be at most the lesser of:
	// (a) the cursor index, or
	// (b) the number of rows minus the maximum number of visible rows (so that
	// as many rows as possible are rendered)
	maximum := max(0, min(m.cursor.n, m.Rows()-m.size))
	m.start = clamp(m.start, minimum, maximum)
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
