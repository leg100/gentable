package gentable

import "github.com/charmbracelet/lipgloss/table"

type window[V any] struct {
	table.Data

	start int
	size  int
}

// At returns the contents of the cell at the given index.
func (m *window[V]) At(row, cell int) string {
	if row >= m.Rows() {
		return ""
	}
	return m.Data.At(row+m.start, cell)
}

func (m *window[V]) Rows() int {
	return min(m.size, m.Data.Rows())
}

func (m *window[V]) PageUp() {
	m.moveStart(-m.size)
}

func (m *window[V]) PageDown() {
	m.moveStart(+m.size)
}

func (m *window[V]) moveStart(n int) {
	if m.size == 0 || m.Data.Rows() == 0 {
		return
	}
	// Move start
	lastRowIndex := m.Data.Rows() - 1
	m.start = clamp(m.start+n, 0, lastRowIndex)
	// Move cursor
	// maxCursor := min(m.start+m.size-1, lastRowIndex)
	// m.cursor.idx = clamp(m.cursor.idx, m.start, maxCursor)
	// m.cursor.id = m.getIDByIndex(m.cursor.idx)
}

//func (m *window[V]) moveCursor(n int) {
//	m.base.moveCursor(n)
//	// Move start
//	if m.size > 0 {
//		startMin := max(0, m.cursor.idx-m.size+1)
//		startMax := min(m.cursor.idx, len(m.rows)-m.size)
//		m.start = clamp(m.start, startMin, startMax)
//	}
//}
