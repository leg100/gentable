package gentable

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	m := booksModel()
	want := strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233003│The German Ideology              │Marx & Engels   │
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233005│Money                            │Martin Amis     │
│390039233006│London Fields                    │Martin Amis     │
│390039233007│Nana                             │Emile Zola      │
│390039233008│To Have and Have Not             │Ernest Hemingway│
│390039233009│The Sun Also Rises               │Ernest Hemingway│
│390039233010│A Farewell to Arms               │Ernest Hemingway│
│390039233011│James Joyce                      │Ulysses         │
│390039233012│Trans-Europe Express             │Owen Hatherley  │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_Sort(t *testing.T) {
	sort := WithSort(func(a, b book) int {
		// sort by author then by their book
		cmp := strings.Compare(a.author, b.author)
		if cmp == 0 {
			return strings.Compare(a.title, b.title)
		}
		return cmp
	})
	m := booksModel(sort)
	want := strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233007│Nana                             │Emile Zola      │
│390039233010│A Farewell to Arms               │Ernest Hemingway│
│390039233009│The Sun Also Rises               │Ernest Hemingway│
│390039233008│To Have and Have Not             │Ernest Hemingway│
│390039233006│London Fields                    │Martin Amis     │
│390039233005│Money                            │Martin Amis     │
│390039233003│The German Ideology              │Marx & Engels   │
│390039233012│Trans-Europe Express             │Owen Hatherley  │
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233011│James Joyce                      │Ulysses         │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_Filter(t *testing.T) {
	m := booksModel()
	m.ApplyFilter(func(b book) bool {
		return b.author == "Ernest Hemingway"
	})
	want := strings.TrimSpace(`
╭────────────┬────────────────────┬────────────────╮
│390039233008│To Have and Have Not│Ernest Hemingway│
│390039233009│The Sun Also Rises  │Ernest Hemingway│
│390039233010│A Farewell to Arms  │Ernest Hemingway│
╰────────────┴────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

	m.RemoveFilter()

	want = strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233003│The German Ideology              │Marx & Engels   │
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233005│Money                            │Martin Amis     │
│390039233006│London Fields                    │Martin Amis     │
│390039233007│Nana                             │Emile Zola      │
│390039233008│To Have and Have Not             │Ernest Hemingway│
│390039233009│The Sun Also Rises               │Ernest Hemingway│
│390039233010│A Farewell to Arms               │Ernest Hemingway│
│390039233011│James Joyce                      │Ulysses         │
│390039233012│Trans-Europe Express             │Owen Hatherley  │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_Height_Before(t *testing.T) {
	m := booksModel()
	m.Height(5)
	want := strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233003│The German Ideology              │Marx & Engels   │
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233005│Money                            │Martin Amis     │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_Height_After(t *testing.T) {
	m := booksModel(WithHeight[book](5))
	want := strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233003│The German Ideology              │Marx & Engels   │
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233005│Money                            │Martin Amis     │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_Height_0(t *testing.T) {
	m := booksModel()
	m.Height(0)
	want := strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233003│The German Ideology              │Marx & Engels   │
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_PageDown(t *testing.T) {
	m := booksModel()
	m.Height(5)

	m.window.PageDown()

	want := strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233005│Money                            │Martin Amis     │
│390039233006│London Fields                    │Martin Amis     │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

	m.window.PageDown()

	want = strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233005│Money                            │Martin Amis     │
│390039233006│London Fields                    │Martin Amis     │
│390039233007│Nana                             │Emile Zola      │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

	m.window.PageDown()

	want = strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233008│To Have and Have Not             │Ernest Hemingway│
│390039233009│The Sun Also Rises               │Ernest Hemingway│
│390039233010│A Farewell to Arms               │Ernest Hemingway│
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

	// no more books so should be a no-op
	m.window.PageDown()

	want = strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233010│A Farewell to Arms               │Ernest Hemingway│
│390039233011│James Joyce                      │Ulysses         │
│390039233012│Trans-Europe Express             │Owen Hatherley  │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_PageUp(t *testing.T) {
	m := booksModel()
	m.Height(5)
	m.window.PageDown()
	m.window.PageDown()
	m.window.PageUp()
	want := strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233005│Money                            │Martin Amis     │
│390039233006│London Fields                    │Martin Amis     │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_BottomTop(t *testing.T) {
	m := booksModel()
	m.Height(5)

	m.window.toBottom()

	want := strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233010│A Farewell to Arms               │Ernest Hemingway│
│390039233011│James Joyce                      │Ulysses         │
│390039233012│Trans-Europe Express             │Owen Hatherley  │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

	m.window.toTop()

	want = strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233003│The German Ideology              │Marx & Engels   │
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233005│Money                            │Martin Amis     │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_Expand(t *testing.T) {
	m := booksModel()
	m.Height(5)

	want := strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233003│The German Ideology              │Marx & Engels   │
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233005│Money                            │Martin Amis     │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

	m.Height(8)

	want = strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233003│The German Ideology              │Marx & Engels   │
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233005│Money                            │Martin Amis     │
│390039233006│London Fields                    │Martin Amis     │
│390039233007│Nana                             │Emile Zola      │
│390039233008│To Have and Have Not             │Ernest Hemingway│
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_Shrink(t *testing.T) {
	m := booksModel()
	m.Height(8)

	want := strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233003│The German Ideology              │Marx & Engels   │
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233005│Money                            │Martin Amis     │
│390039233006│London Fields                    │Martin Amis     │
│390039233007│Nana                             │Emile Zola      │
│390039233008│To Have and Have Not             │Ernest Hemingway│
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

	m.Height(5)

	want = strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬────────────────╮
│390039233003│The German Ideology              │Marx & Engels   │
│390039233004│Death in Venice and Other Stories│Thomas Mann     │
│390039233005│Money                            │Martin Amis     │
╰────────────┴─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

}

func booksModel(opts ...Option[book]) Model[book] {
	m := New[book]()
	m.Border(lipgloss.RoundedBorder())
	for _, fn := range opts {
		fn(&m)
	}
	for _, bk := range books {
		m.Append(Row[book]{
			V: bk,
			Cells: []string{
				string(bk.isbn),
				bk.title,
				bk.author,
			},
		})
	}
	return m
}
