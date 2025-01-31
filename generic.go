package gentable

import "github.com/charmbracelet/lipgloss/table"

type ID any

type Generic[V any] struct {
	data   *table.StringData
	render func(V) []string
	sort   func(V, V) int
}

func NewGeneric[V any](
	render func(V) []string,
) *Generic[V] {
	return &Generic[V]{
		data:   table.NewStringData(),
		render: render,
	}
}

func (g *Generic[V]) Append(rows ...V) {
	for _, row := range rows {
		g.data.Append(g.render(row))
	}
}

func (g *Generic[V]) WithSort(fn func(a, b V) int) {
	g.sort = fn
}

func (g *Generic[V]) At(row, cell int) string {
	return g.data.At(row, cell)
}

func (g *Generic[V]) Columns() int {
	return g.data.Columns()
}

func (g *Generic[V]) Rows() int {
	return g.data.Rows()
}
