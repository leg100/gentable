package gentable

import (
	"slices"
)

type ID any

type base[V any] struct {
	rows    []ID
	db      map[ID]V
	cells   map[ID][]string
	columns int
	getID   func(V) ID
	render  func(V) []string
	sort    func(V, V) int
}

func newBase[V any](
	getID func(V) ID,
	render func(V) []string,
) *base[V] {
	return &base[V]{
		db:     make(map[ID]V),
		cells:  make(map[ID][]string),
		getID:  getID,
		render: render,
	}
}

func (m *base[V]) Append(rows ...V) {
	for _, row := range rows {
		cells := m.render(row)
		m.columns = max(m.columns, len(cells))
		newID := m.getID(row)
		m.rows = append(m.rows, newID)
		m.db[newID] = row
		m.cells[newID] = cells
	}
	//if m.cursor.id == nil {
	//	m.cursor.id = m.rows[0]
	//}
	if m.sort != nil {
		slices.SortFunc(m.rows, func(a, b ID) int {
			return m.sort(m.db[a], m.db[b])
		})
		//for i, id := range m.rows {
		//	if id == m.cursor.id {
		//		m.moveCursor(i - m.cursor.idx)
		//		break
		//	}
		//}
	}
}

// At returns the contents of the cell at the given index.
func (m *base[V]) At(row, cell int) string {
	if row >= len(m.rows) {
		return ""
	}
	id := m.rows[row]
	if cell >= len(m.cells[id]) {
		return ""
	}
	return m.cells[id][cell]
}

func (m *base[V]) Columns() int {
	return m.columns
}

func (m *base[V]) Rows() int {
	return len(m.rows)
}

//func (m *base[V]) moveCursor(n int) {
//	if n == 0 || len(m.rows) == 0 {
//		return
//	}
//	// Move cursor
//	lastRowIndex := len(m.rows) - 1
//	m.cursor.idx = clamp(m.cursor.idx+n, 0, lastRowIndex)
//	m.cursor.id = m.rows[m.cursor.idx]
//}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
