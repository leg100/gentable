package gentable

import (
	"testing"

	"github.com/charmbracelet/x/exp/golden"
)

func TestModel(t *testing.T) {
	m := newTestModel(t)

	golden.RequireEqual(t, []byte(m.View()))
}

func TestModel_withScrollbar(t *testing.T) {
	m := newTestModel(t, WithScrollbar[tx]())

	golden.RequireEqual(t, []byte(m.View()))
}

func TestModel_narrow(t *testing.T) {
	m := newTestModel(t)
	m.SetDimensions(30, 10)

	golden.RequireEqual(t, []byte(m.View()))
}
