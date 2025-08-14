package gentable

import "slices"

// base is the complete set of rows
type base[V comparable] struct {
	rows    []V
	cells   map[V][]string
	columns int
	sort    func(V, V) int
}

type row[V comparable] struct {
	v     V
	cells []string
}

func (m *base[V]) Append(rows ...row[V]) {
	for _, row := range rows {
		m.cells[row.v] = row.cells
		m.columns = max(m.columns, len(row.cells))
		m.rows = append(m.rows, row.v)

		//if m.cursor.id == nil {
		//	m.cursor.id = m.rows[0]
		//}
	}
	if m.sort != nil {
		slices.SortFunc(m.rows, m.sort)
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
	return m.getCell(m.rows[row], cell)
}

func (m *base[V]) getCell(id ID, cell int) string {
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
