package gentable

// base is the complete set of rows
type base[V comparable] struct {
	rows    []V
	cells   map[V][]string
	columns int
}

type row[V comparable] struct {
	v     V
	cells []string
}

func (m *base[V]) Append(rows ...row[V]) {
	for _, row := range rows {
		m.rows = append(m.rows, row.v)
		m.cells[row.v] = row.cells
		m.columns = max(m.columns, len(row.cells))
	}
}

func (m *base[V]) Columns() int {
	return m.columns
}

func (m *base[V]) Rows() []V {
	return m.rows
}

func (m *base[V]) getCells(v V) []string {
	return m.cells[v]
}
