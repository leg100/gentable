package gentable

import "github.com/charmbracelet/lipgloss/table"

type data struct {
	*table.StringData
	// index of first row in current window.
	start int
	// cursor is the index of the current row (typically highlighted in the UI
	// for the user's benefit).
	cursor int
	// size of the "window" onto the underlying data. If negative then the size
	// of the underlying data is used.
	size int
}

func newData() *data {
	return &data{StringData: table.NewStringData(), size: -1}
}

func (d *data) Rows() int {
	if d.size < 0 {
		return d.StringData.Rows()
	}
	return min(d.size, d.StringData.Rows())
}

// At returns the contents of the cell at the given index.
func (m *data) At(row, cell int) string {
	if m.size < 0 {
		return m.StringData.At(row, cell)
	}
	row += m.start
	return m.StringData.At(row, cell)
}

//	func (m *data) MoveCursor(n int) {
//		if n == 0 {
//			return
//		}
//		if n.
//			return
//		}
//		m.start = min(m.start+m.size, max(0, m.StringData.Rows()-1))
//	}
func (m *data) PageUp() {
	if m.size < 0 {
		return
	}
	m.start = max(0, m.start-m.size)
}

func (m *data) PageDown() {
	if m.size < 0 {
		return
	}
	if m.start > max(0, m.StringData.Rows()-1) {
		return
	}
	m.start = min(m.start+m.size, max(0, m.StringData.Rows()-1))
}
