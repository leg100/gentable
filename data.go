package gentable

import "github.com/charmbracelet/lipgloss/table"

type data[V any] struct {
	table.Data
}

func NewData[V any](
	getID func(V) ID,
	render func(V) []string,
) *data[V] {
	return &data[V]{
		Data: newBase(getID, render),
	}
}
