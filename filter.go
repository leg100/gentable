package gentable

import "slices"

type filter[V any] struct {
	base[V]
	filtered []ID
	filterfn func(V) bool
}

func newFilter[V any](base base[V], fn func(V) bool) *filter[V] {
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
	for _, v := range rows {
		if m.filterfn(v) {
			m.filtered = append(m.filtered, v)
		}
	}
	if m.sort != nil {
		slices.SortFunc(m.filtered, func(a, b ID) int {
			return m.sort(m.db[a], m.db[b])
		})
	}
}

func (m *filter[V]) Rows() int {
	return len(m.filtered)
}
