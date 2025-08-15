package gentable

import (
	"strings"
	"testing"

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

func TestModel_withSort(t *testing.T) {
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

func TestModel_Height(t *testing.T) {
	m := booksModel()
	m.Height(5)
	want := strings.TrimSpace(`
╭────────────┬─────────────────────────────────┬─────────────╮
│390039233003│The German Ideology              │Marx & Engels│
│390039233004│Death in Venice and Other Stories│Thomas Mann  │
│390039233005│Money                            │Martin Amis  │
╰────────────┴─────────────────────────────────┴─────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_Height_0(t *testing.T) {
	m := booksModel()
	m.Height(0)
	want := strings.TrimSpace(`
╭────────────┬───────────────────┬─────────────╮
│390039233003│The German Ideology│Marx & Engels│
╰────────────┴───────────────────┴─────────────╯
`)
	assert.Equal(t, want, m.View())
}

//func TestModel_PageUp(t *testing.T) {
//	m := modelWithWindowStringData()
//	m.Height(5)
//	m.data.(*window).PageDown()
//	want := strings.TrimSpace(`
//╭────────────────────┬────────────────╮
//│London Fields       │Martin Amis     │
//│Nana                │Emile Zola      │
//│To Have and Have Not│Ernest Hemingway│
//╰────────────────────┴────────────────╯
//`)
//	assert.Equal(t, want, m.View())
//	assert.Equal(t, 3, m.data.(*window).cursor)
//
//	m.data.(*window).PageDown()
//	want = strings.TrimSpace(`
//╭──────────────────┬────────────────╮
//│The Sun Also Rises│Ernest Hemingway│
//│A Farewell to Arms│Ernest Hemingway│
//│James Joyce       │Ulysses         │
//╰──────────────────┴────────────────╯
//`)
//	assert.Equal(t, want, m.View())
//	assert.Equal(t, 6, m.data.(*window).cursor)
//
//	m.data.(*window).PageDown()
//	want = strings.TrimSpace(`
//╭────────────────────┬──────────────╮
//│Trans-Europe Express│Owen Hatherley│
//│                    │              │
//│                    │              │
//╰────────────────────┴──────────────╯
//`)
//	assert.Equal(t, want, m.View())
//	assert.Equal(t, 9, m.data.(*window).cursor)
//
//	m.data.(*window).PageDown()
//	want = strings.TrimSpace(`
//╭────────────────────┬──────────────╮
//│Trans-Europe Express│Owen Hatherley│
//│                    │              │
//│                    │              │
//╰────────────────────┴──────────────╯
//`)
//	assert.Equal(t, want, m.View())
//	assert.Equal(t, 9, m.data.(*window).cursor)
//
//	m.data.(*window).PageUp()
//	want = strings.TrimSpace(`
//╭──────────────────┬────────────────╮
//│The Sun Also Rises│Ernest Hemingway│
//│A Farewell to Arms│Ernest Hemingway│
//│James Joyce       │Ulysses         │
//╰──────────────────┴────────────────╯
//`)
//	assert.Equal(t, want, m.View())
//	assert.Equal(t, 8, m.data.(*window).cursor)
//
//	m.data.(*window).PageUp()
//	want = strings.TrimSpace(`
//╭────────────────────┬────────────────╮
//│London Fields       │Martin Amis     │
//│Nana                │Emile Zola      │
//│To Have and Have Not│Ernest Hemingway│
//╰────────────────────┴────────────────╯
//`)
//	assert.Equal(t, want, m.View())
//	assert.Equal(t, 5, m.data.(*window).cursor)
//
//	m.data.(*window).PageUp()
//	want = strings.TrimSpace(`
//╭─────────────────────────────────┬─────────────╮
//│The German Ideology              │Marx & Engels│
//│Death in Venice and Other Stories│Thomas Mann  │
//│Money                            │Martin Amis  │
//╰─────────────────────────────────┴─────────────╯
//`)
//	assert.Equal(t, want, m.View())
//	assert.Equal(t, 2, m.data.(*window).cursor)
//
//	m.data.(*window).PageUp()
//	want = strings.TrimSpace(`
//╭─────────────────────────────────┬─────────────╮
//│The German Ideology              │Marx & Engels│
//│Death in Venice and Other Stories│Thomas Mann  │
//│Money                            │Martin Amis  │
//╰─────────────────────────────────┴─────────────╯
//`)
//	assert.Equal(t, want, m.View())
//	assert.Equal(t, 2, m.data.(*window).cursor)
//}

func booksModel(opts ...Option[book]) Model[book] {
	m := New[book]()
	for _, fn := range opts {
		fn(&m)
	}
	for _, bk := range books {
		m.Append(row[book]{
			v: bk,
			cells: []string{
				string(bk.isbn),
				bk.title,
				bk.author,
			},
		})
	}
	return m
}

//func modelWithWindowStringData() Model {
//	data := newData(
//		func(v []string) ID {
//			id := i
//			i++
//			return id
//		},
//		func(v []string) []string { return v },
//	)
//	for _, book := range books {
//		data.Append([]string{
//			book.title,
//			book.author,
//		})
//	}
//	win := window{Data: data}
//	return New(&win)
//}
