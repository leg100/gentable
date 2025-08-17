package gentable

type data[V comparable] struct {
	unfiltered []V

	cells   map[V][]string
	columns int

	filtered []V
	filter   func(V) bool
}

type Row[V comparable] struct {
	V     V
	Cells []string
}

func (d *data[V]) At(row, cell int) string {
	v := d.rows()[row]
	cells := d.cells[v]
	if cell >= len(cells) {
		// Not all rows have the same number of cells but the lipgloss table lib
		// that calls this method doesn't know that.
		return ""
	}
	return cells[cell]
}

func (d *data[V]) Rows() int { return len(d.rows()) }

func (b *data[V]) Columns() int { return b.columns }

func (b *data[V]) add(rows ...Row[V]) {
	for _, row := range rows {
		b.unfiltered = append(b.unfiltered, row.V)
		b.cells[row.V] = row.Cells
		b.columns = max(b.columns, len(row.Cells))

		if b.filter != nil {
			if b.filter(row.V) {
				b.filtered = append(b.filtered, row.V)
			}
		}
	}
}

func (b *data[V]) rows() []V {
	if b.filter != nil {
		return b.filtered
	}
	return b.unfiltered
}

func (b *data[V]) applyFilter(fn func(V) bool) {
	b.filter = fn
	for _, v := range b.unfiltered {
		if fn(v) {
			b.filtered = append(b.filtered, v)
		}
	}
}

func (b *data[V]) removeFilter() {
	b.filter = nil
	b.filtered = nil
}
