package gentable

import "github.com/charmbracelet/lipgloss/table"

type data[V any] struct {
	table.Data
	filterfn func(V) bool
}

func NewData[V any](
	getID func(V) ID,
	render func(V) []string,
) *data[V] {
	return &data[V]{
		Data: newBase(getID, render),
	}
}

func (d *data[V]) toggleFilter() {
	switch data := d.Data.(type) {
	case *filter[V]:
		d.Data = data.base
	case *base[V]:
		d.Data = newFilter(data, d.filterfn)
	}
}
