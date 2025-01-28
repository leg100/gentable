package gentable

type Option[V any] func(m *Model[V])

// WithSortFunc configures the table to sort rows using the given func.
func WithSortFunc[V any](sortFunc func(V, V) int) Option[V] {
	return func(m *Model[V]) {
		m.sortFunc = sortFunc
	}
}

// WithSelectable sets whether rows are selectable.
func WithSelectable[V any](s bool) Option[V] {
	return func(m *Model[V]) {
		m.selectable = s
	}
}

// WithColumns defines columns. The order given is the order in which columns
// will be rendered from left to right.
func WithColumns[V any](cols ...Column) Option[V] {
	return func(m *Model[V]) {
		// Copy column structs onto receiver, because the caller may modify
		// columns afterwards.
		m.cols = make([]Column, len(cols))
		copy(m.cols, cols)
		// For each column, set default truncation function if unset.
		for i, col := range m.cols {
			if col.TruncationFunc == nil {
				m.cols[i].TruncationFunc = defaultTruncationFunc
			}
		}
	}
}

// WithScrollbar adds a scrollbar to the right hand side.
func WithScrollbar[V any]() Option[V] {
	return func(m *Model[V]) {
		m.scrollbar = true
	}
}

// WithBorder adds a border around the table.
func WithBorder[V any]() Option[V] {
	return func(m *Model[V]) {
		m.border = true
	}
}
