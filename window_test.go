package gentable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWindow(t *testing.T) {
	type isbn int

	t.Run("set cursor to first added row", func(t *testing.T) {
		win := newWindow[isbn](0)
		win.Append(Row[isbn]{v: 123})
		win.Append(Row[isbn]{v: 456})

		assert.Equal(t, 0, win.cursor.n)
		assert.Equal(t, isbn(123), win.cursor.v)
	})

	t.Run("move cursor down one row", func(t *testing.T) {
		win := newWindow[isbn](0)
		win.Append(Row[isbn]{v: 123})
		win.Append(Row[isbn]{v: 456})

		win.moveCursor(1)

		assert.Equal(t, 1, win.cursor.n)
		assert.Equal(t, isbn(456), win.cursor.v)
	})

	t.Run("move cursor down and up one row", func(t *testing.T) {
		win := newWindow[isbn](0)
		win.Append(Row[isbn]{v: 123})
		win.Append(Row[isbn]{v: 456})

		win.moveCursor(1)
		win.moveCursor(-1)

		assert.Equal(t, 0, win.cursor.n)
		assert.Equal(t, isbn(123), win.cursor.v)
	})

	t.Run("cursor cannot move beyond last row", func(t *testing.T) {
		win := newWindow[isbn](0)
		win.Append(Row[isbn]{v: 123})
		win.Append(Row[isbn]{v: 456})
		win.Append(Row[isbn]{v: 789})

		win.moveCursor(999)

		assert.Equal(t, 2, win.cursor.n)
		assert.Equal(t, isbn(789), win.cursor.v)
	})

	t.Run("cursor cannot move above first row", func(t *testing.T) {
		win := newWindow[isbn](0)
		win.Append(Row[isbn]{v: 123})

		win.moveCursor(-999)

		assert.Equal(t, 0, win.cursor.n)
		assert.Equal(t, isbn(123), win.cursor.v)
	})

	t.Run("cursor cannot move above first row", func(t *testing.T) {
		win := newWindow[isbn](0)
		win.Append(Row[isbn]{v: 123})

		win.moveCursor(-999)

		assert.Equal(t, 0, win.cursor.n)
		assert.Equal(t, isbn(123), win.cursor.v)
	})

	t.Run("moving cursor beyond window moves window down", func(t *testing.T) {
		win := newWindow[isbn](3)
		win.Append(Row[isbn]{v: 12})
		win.Append(Row[isbn]{v: 34})
		win.Append(Row[isbn]{v: 56})
		win.Append(Row[isbn]{v: 78})
		win.Append(Row[isbn]{v: 90})

		win.moveCursor(3)

		assert.Equal(t, 1, win.start)
		assert.Equal(t, 3, win.cursor.n)
		assert.Equal(t, isbn(78), win.cursor.v)
	})

	t.Run("moving cursor above window moves window up", func(t *testing.T) {
		win := newWindow[isbn](3)
		win.Append(Row[isbn]{v: 12})
		win.Append(Row[isbn]{v: 34})
		win.Append(Row[isbn]{v: 56})
		win.Append(Row[isbn]{v: 78})
		win.Append(Row[isbn]{v: 90})

		// move cursor beyond window to move window down
		win.moveCursor(4)

		// first visible row is now index 2
		assert.Equal(t, 2, win.start)

		win.moveCursor(-3)

		// first visible row is now index 1
		assert.Equal(t, 1, win.start)
	})
}
