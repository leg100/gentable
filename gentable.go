package gentable

import "github.com/charmbracelet/lipgloss/table"

type data interface {
	table.Data
	toggleFilter()
}

//func (d *db[V]) toggleFilter() {
//	switch data := d.Data.(type) {
//	case *filter[V]:
//		d.Data = data.base
//	case *base[V]:
//		d.Data = newFilter(data, d.filterfn)
//	}
//}
