package gentable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestData_cursor(t *testing.T) {
	tests := []struct {
		name string
		do   func(*db[book])
		want func(*db[book])
	}{
		{
			name: "set cursor to first added row",
			do: func(d *db[book]) {
				d.Append(
					book{author: "marx", title: "capital", isbn: isbn("123")},
					book{author: "joyce", title: "dubliners", isbn: isbn("456")},
				)
			},
			want: func(d *db[book]) {
				assert.Equal(t, d.cursor.idx, 0)
				assert.Equal(t, d.cursor.id, isbn("123"))
			},
		},
		{
			name: "move cursor down one row",
			do: func(d *db[book]) {
				d.Append(
					books[0],
					books[1],
				)
				d.moveCursor(1)
			},
			want: func(d *db[book]) {
				assert.Equal(t, d.cursor.idx, 1)
				assert.Equal(t, d.cursor.id, isbn("390039233004"))
			},
		},
		{
			name: "move cursor down and up one row",
			do: func(d *db[book]) {
				d.Append(
					books[0],
					books[1],
				)
				d.moveCursor(1)
				d.moveCursor(-1)
			},
			want: func(d *db[book]) {
				assert.Equal(t, d.cursor.idx, 0)
				assert.Equal(t, d.cursor.id, isbn("390039233003"))
			},
		},
		{
			name: "cursor cannot move beyond last row",
			do: func(d *db[book]) {
				d.Append(books[0])
				d.moveCursor(1)
			},
			want: func(d *db[book]) {
				assert.Equal(t, d.cursor.idx, 0)
				assert.Equal(t, d.cursor.id, isbn("390039233003"))
			},
		},
		{
			name: "cursor cannot move before first row",
			do: func(d *db[book]) {
				d.Append(books[0])
				d.moveCursor(-1)
			},
			want: func(d *db[book]) {
				assert.Equal(t, d.cursor.idx, 0)
				assert.Equal(t, d.cursor.id, isbn("390039233003"))
			},
		},
		{
			name: "move window beyond cursor moves cursor",
			do: func(d *db[book]) {
				d.SetWindowSize(1)
				d.Append(
					books[0],
					books[1],
				)
				d.moveStart(1)
			},
			want: func(d *db[book]) {
				assert.Equal(t, d.cursor.idx, 1)
				assert.Equal(t, d.cursor.id, isbn("390039233004"))
			},
		},
		{
			name: "move window but cursor does not need moving",
			do: func(d *db[book]) {
				d.SetWindowSize(2)
				d.Append(
					books[0],
					books[1],
					books[2],
				)
				d.moveCursor(1)
				d.moveStart(1)
			},
			want: func(d *db[book]) {
				assert.Equal(t, d.cursor.idx, 1)
				assert.Equal(t, d.cursor.id, isbn("390039233004"))
			},
		},
		{
			name: "setting window restricts visible rows and columns",
			do: func(d *db[book]) {
				d.SetWindowSize(2)
				d.Append(
					book{author: "marx", title: "capital", isbn: isbn("123")},
					book{author: "joyce", title: "dubliners", isbn: isbn("456")},
					book{author: "bukowski", title: "ham & rye", isbn: isbn("789")},
				)
			},
			want: func(d *db[book]) {
				assert.Equal(t, 2, d.Rows())
				assert.Equal(t, "marx", d.At(0, 0))
				assert.Equal(t, "joyce", d.At(1, 0))
				assert.Equal(t, "", d.At(2, 0))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := newData(
				getISBN,
				renderBook,
			)
			tt.do(data)
			tt.want(data)
		})
	}
}
