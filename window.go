package gentable

import "github.com/charmbracelet/lipgloss/table"

type window struct {
	table.Data

	// index of first row in window.
	start int
	// cursor is the index of the current row.
	cursor int
	// size of the "window" onto the underlying data. If zero then the size
	// of the underlying data is used.
	size int
}

func (d *window) SetWindowSize(size int) {
	d.size = size
}

func (d *window) Rows() int {
	if d.size == 0 {
		return d.Data.Rows()
	}
	return min(d.size, d.Data.Rows())
}

// At returns the contents of the cell at the given index.
func (m *window) At(row, cell int) string {
	if m.size == 0 {
		return m.Data.At(row, cell)
	}
	row += m.start
	return m.Data.At(row, cell)
}

func (m *window) PageUp() {
	m.moveStart(-m.size)
}

func (m *window) PageDown() {
	m.moveStart(+m.size)
}

func (m *window) moveStart(n int) {
	if m.size == 0 || m.Data.Rows() == 0 {
		return
	}

	// Move start
	lastRowIndex := m.Data.Rows() - 1
	m.start = clamp(m.start+n, 0, lastRowIndex)

	// Move cursor
	maxCursor := min(m.start+m.size-1, lastRowIndex)
	m.cursor = clamp(m.cursor, m.start, maxCursor)
}

func (m *window) moveCursor(n int) {
	if n == 0 || m.Data.Rows() == 0 {
		return
	}

	// Move cursor
	lastRowIndex := m.Data.Rows() - 1
	m.cursor = clamp(m.cursor+n, 0, lastRowIndex)

	// Move start
	startMin := max(0, m.cursor-m.size+1)
	startMax := min(m.cursor, m.Data.Rows()-m.size)
	m.start = clamp(m.start, startMin, startMax)
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
