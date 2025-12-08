package gentable

// data is the database for a table.
type data[V comparable] struct {
	unfiltered []V

	cells   map[V][]string
	columns int

	filtered []V
	filter   func(V) bool

	selected map[V]bool

	headers []string

	*resizer
}

type Row[V comparable] struct {
	V     V
	Cells []string
}

func newData[V comparable](headers ...string) *data[V] {
	return &data[V]{
		cells:    make(map[V][]string),
		selected: make(map[V]bool),
		resizer:  newResizer(headers...),
		// keep record of headers for when resizer needs to be re-initialized.
		headers: headers,
	}
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
	cells := make([][]string, len(rows))

	for i, row := range rows {
		b.unfiltered = append(b.unfiltered, row.V)
		b.cells[row.V] = row.Cells
		b.columns = max(b.columns, len(row.Cells))

		if b.filter != nil {
			if b.filter(row.V) {
				b.filtered = append(b.filtered, row.V)
			}
		}

		cells[i] = row.Cells
	}

	b.resize(false, cells...)
}

func (b *data[V]) rows() []V {
	if b.filter != nil {
		return b.filtered
	}
	return b.unfiltered
}

func (b *data[V]) applyFilter(fn func(V) bool) {
	b.filter = fn

	// Track selections that pass the filter
	selected := make(map[V]bool)

	for _, v := range b.unfiltered {
		if fn(v) {
			b.filtered = append(b.filtered, v)
			selected[v] = b.selected[v]
		}
	}

	b.selected = selected
}

func (b *data[V]) removeFilter() {
	b.filter = nil
	b.filtered = nil
}

func (b *data[V]) toggleSelect(row int) {
	v := b.rows()[row]
	if b.selected[v] {
		delete(b.selected, v)
	} else {
		b.selected[v] = true
	}
}

func (b *data[V]) selectAll() {
	for _, row := range b.rows() {
		b.selected[row] = true
	}
}

func (b *data[V]) clearSelection() {
	b.selected = make(map[V]bool)
}

func (b *data[V]) isSelected(row int) bool {
	// The lipgloss stylefunc which calls this method can send negative row
	// number (indicating the header row).
	if row < 0 || row > len(b.rows()) {
		return false
	}
	v := b.rows()[row]
	return b.selected[v]
}
