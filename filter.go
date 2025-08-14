package gentable

type filter[V comparable] struct {
	data[V]

	filtered []V
	filter   func(V) bool
}

func newFilter[V comparable](db data[V], fn func(V) bool) *filter[V] {
	f := &filter[V]{
		data:   db,
		filter: fn,
	}
	for _, v := range db.Rows() {
		if fn(v) {
			f.filtered = append(f.filtered, v)
		}
	}
	return f
}

func (m *filter[V]) Append(rows ...row[V]) {
	m.data.Append(rows...)
	for _, row := range rows {
		if m.filter(row.v) {
			m.filtered = append(m.filtered, row.v)
		}
	}
}

func (m *filter[V]) Rows() []V {
	return m.filtered
}
