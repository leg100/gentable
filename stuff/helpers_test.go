package gentable

import (
	"strconv"
	"testing"
	"time"
)

func newTestModel(t *testing.T, opts ...Option[tx]) Model[tx] {
	opts = append(opts,
		WithColumns[tx](
			Column{
				Key:   "price",
				Title: "PRICE",
			},
			Column{
				Key:   "date",
				Title: "DATE",
			},
			Column{
				Key:   "post_code",
				Title: "POST CODE",
			},
		),
	)
	m := New(
		func(p tx) ID {
			return p.id
		},
		func(p tx) RenderedCells {
			return map[ColumnKey]string{
				"price":     strconv.Itoa(p.pounds),
				"date":      p.date.Format(time.DateOnly),
				"post_code": p.postCode,
			}
		},
		40,
		10,
		opts...,
	)
	m.AddItems(txs...)
	return m
}
