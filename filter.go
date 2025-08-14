package gentable

type filter[V comparable] struct {
	db       filterDB[V]
	filtered []ID
	filterfn func(V) bool
}

type filterDB[V comparable] interface {
	getValue(ID) V
	getCells(ID) []string
	sort([]ID)
}

func newFilter[V any](db filterDB[V], ids []ID, fn func(V) bool) *filter[V] {
	filter := &filter[V]{
		db:       db,
		filterfn: fn,
	}
	filter.append(ids...)
	return filter
}

func (m *filter[V]) Append(rows ...row[V]) {
	for _, id := range ids {
		if m.filterfn(m.db.getValue(id)) {
			m.filtered = append(m.filtered, id)
		}
	}
	m.db.sort(m.filtered)
}

// At returns the contents of the cell at the given index.
func (m *filter[V]) At(row, cell int) string {
	if row >= len(m.filtered) {
		return ""
	}
	return m.db.getCells(m.filtered[row])[cell]
}

func (m *filter[V]) Rows() int {
	return len(m.filtered)
}
