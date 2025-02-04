package gentable

import "slices"

type filter[V any] struct {
	*base[V]
	filtered []ID
	filterfn func(V) bool
	columns  int
}

func newFilter[V any](base *base[V], fn func(V) bool) *filter[V] {
	var filtered []ID
	for _, v := range base.db {
		if fn(v) {
			filtered = append(filtered, v)
		}
	}
	return &filter[V]{
		base:     base,
		filtered: filtered,
		filterfn: fn,
	}
}

func (m *filter[V]) Append(rows ...V) {
	m.base.Append(rows...)
	for _, id := range rows {
		v := m.db[id]
		if m.filterfn(v) {
			m.filtered = append(m.filtered, id)
		}
		m.columns = max(m.columns, len(m.cells[id]))
	}
	if m.sort != nil {
		slices.SortFunc(m.filtered, func(a, b ID) int {
			return m.sort(m.db[a], m.db[b])
		})
	}
}

// At returns the contents of the cell at the given index.
func (m *filter[V]) At(row, cell int) string {
	if row >= len(m.filtered) {
		return ""
	}
	return m.getCell(m.filtered[row], cell)
}

func (m *filter[V]) Rows() int {
	return len(m.filtered)
}

func (m *filter[V]) Columns() int {
	return m.columns
}
