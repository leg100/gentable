package gentable

type Option[V comparable] func(m *Model[V])

func WithSort[V comparable](fn func(V, V) int) Option[V] {
	return func(m *Model[V]) {
		m.sort = fn
	}
}

func WithHeight[V comparable](height int) Option[V] {
	return func(m *Model[V]) {
		m.Height(height)
	}
}
