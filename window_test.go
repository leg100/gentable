package gentable

import (
	"testing"

	"github.com/charmbracelet/lipgloss/v2/table"
	"github.com/stretchr/testify/assert"
)

func TestWindow(t *testing.T) {
	type isbn int

	t.Run("set cursor to first added row", func(t *testing.T) {
		win := newWindow[isbn](&testWindowTable{})
		win.Append(Row[isbn]{V: 123})
		win.Append(Row[isbn]{V: 456})

		assert.Equal(t, 0, win.cursor.n)
		assert.Equal(t, isbn(123), win.cursor.v)
	})

	t.Run("move cursor down one row", func(t *testing.T) {
		win := newWindow[isbn](&testWindowTable{})
		win.Append(Row[isbn]{V: 123})
		win.Append(Row[isbn]{V: 456})

		win.moveCursor(1)

		assert.Equal(t, 1, win.cursor.n)
		assert.Equal(t, isbn(456), win.cursor.v)
	})

	t.Run("move cursor down and up one row", func(t *testing.T) {
		win := newWindow[isbn](&testWindowTable{})
		win.Append(Row[isbn]{V: 123})
		win.Append(Row[isbn]{V: 456})

		win.moveCursor(1)
		win.moveCursor(-1)

		assert.Equal(t, 0, win.cursor.n)
		assert.Equal(t, isbn(123), win.cursor.v)
	})

	t.Run("cursor cannot move beyond last row", func(t *testing.T) {
		win := newWindow[isbn](&testWindowTable{})
		win.Append(Row[isbn]{V: 123})
		win.Append(Row[isbn]{V: 456})
		win.Append(Row[isbn]{V: 789})

		win.moveCursor(999)

		assert.Equal(t, 2, win.cursor.n)
		assert.Equal(t, isbn(789), win.cursor.v)
	})

	t.Run("cursor cannot move above first row", func(t *testing.T) {
		win := newWindow[isbn](&testWindowTable{})
		win.Append(Row[isbn]{V: 123})

		win.moveCursor(-999)

		assert.Equal(t, 0, win.cursor.n)
		assert.Equal(t, isbn(123), win.cursor.v)
	})

	t.Run("cursor cannot move above first row", func(t *testing.T) {
		win := newWindow[isbn](&testWindowTable{})
		win.Append(Row[isbn]{V: 123})

		win.moveCursor(-999)

		assert.Equal(t, 0, win.cursor.n)
		assert.Equal(t, isbn(123), win.cursor.v)
	})

	t.Run("moving cursor beyond window moves window down", func(t *testing.T) {
		win := newWindow[isbn](&testWindowTable{size: 3})
		win.Append(Row[isbn]{V: 12})
		win.Append(Row[isbn]{V: 34})
		win.Append(Row[isbn]{V: 56})
		win.Append(Row[isbn]{V: 78})
		win.Append(Row[isbn]{V: 90})

		win.moveCursor(3)

		assert.Equal(t, 1, win.GetYOffset())
		assert.Equal(t, 3, win.cursor.n)
		assert.Equal(t, isbn(78), win.cursor.v)
	})

	t.Run("moving cursor above window moves window up", func(t *testing.T) {
		win := newWindow[isbn](&testWindowTable{size: 3})
		win.Append(Row[isbn]{V: 12})
		win.Append(Row[isbn]{V: 34})
		win.Append(Row[isbn]{V: 56})
		win.Append(Row[isbn]{V: 78})
		win.Append(Row[isbn]{V: 90})

		// move cursor beyond window to move window down
		win.moveCursor(4)

		// first visible row is now index 2
		assert.Equal(t, 2, win.GetYOffset())

		win.moveCursor(-3)

		// first visible row is now index 1
		assert.Equal(t, 1, win.GetYOffset())
	})

	t.Run("shrinking window moves cursor", func(t *testing.T) {
		// window with all rows visible
		table := &testWindowTable{size: 5}
		win := newWindow[isbn](table)
		win.Append(Row[isbn]{V: 12})
		win.Append(Row[isbn]{V: 34})
		win.Append(Row[isbn]{V: 56})
		win.Append(Row[isbn]{V: 78})
		win.Append(Row[isbn]{V: 90})

		// move cursor to last row
		win.moveCursor(4)

		// shrink window to 3 rows
		table.size = 3
		win.reset()

		// first visible row should now be index 2
		assert.Equal(t, 2, win.GetYOffset())

		win.moveCursor(-3)

		// first visible row is now index 1
		assert.Equal(t, 1, win.GetYOffset())
	})
}

type testWindowTable struct {
	yoffset int
	size    int
}

func (t *testWindowTable) GetYOffset() int {
	return t.yoffset
}

func (t *testWindowTable) YOffset(o int) *table.Table {
	t.yoffset = o
	return nil
}

func (t *testWindowTable) VisibleRows() int {
	return t.size
}
