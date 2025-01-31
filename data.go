package gentable

import "github.com/charmbracelet/lipgloss/table"

type data struct {
	*table.StringData
	// index of first row in window.
	start int
	// cursor is the index of the current row.
	cursor int
	// size of the "window" onto the underlying data. If zero then the size
	// of the underlying data is used.
	size int
}

func newData() *data {
	return &data{StringData: table.NewStringData()}
}

func (d *data) Rows() int {
	if d.size == 0 {
		return d.StringData.Rows()
	}
	return min(d.size, d.StringData.Rows())
}

// At returns the contents of the cell at the given index.
func (m *data) At(row, cell int) string {
	if m.size == 0 {
		return m.StringData.At(row, cell)
	}
	row += m.start
	return m.StringData.At(row, cell)
}

func (m *data) moveCursor(n int) {
	if n == 0 || m.size == 0 || m.StringData.Rows() == 0 {
		return
	}

	// Move cursor
	lastRowIndex := m.StringData.Rows() - 1
	m.cursor = clamp(m.cursor+n, 0, lastRowIndex)

	// Move start
	startMin := max(0, m.cursor-m.size+1)
	startMax := min(m.cursor, m.StringData.Rows()-m.size)
	m.start = clamp(m.start, startMin, startMax)
}

func (m *data) PageUp() {
	m.moveStart(-m.size)
}

func (m *data) PageDown() {
	m.moveStart(+m.size)
}

func (m *data) moveStart(n int) {
	if m.size == 0 || m.StringData.Rows() == 0 {
		return
	}

	// Move start
	lastRowIndex := m.StringData.Rows() - 1
	m.start = clamp(m.start+n, 0, lastRowIndex)

	// Move cursor
	maxCursor := min(m.start+m.size-1, lastRowIndex)
	m.cursor = clamp(m.cursor, m.start, maxCursor)
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
