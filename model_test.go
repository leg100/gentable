package gentable

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	m := booksModel()
	want := strings.TrimSpace(`
╭─────────────────────────────────┬────────────────╮
│The German Ideology              │Marx & Engels   │
│Death in Venice and Other Stories│Thomas Mann     │
│Money                            │Martin Amis     │
│London Fields                    │Martin Amis     │
│Nana                             │Emile Zola      │
│To Have and Have Not             │Ernest Hemingway│
│The Sun Also Rises               │Ernest Hemingway│
│A Farewell to Arms               │Ernest Hemingway│
│James Joyce                      │Ulysses         │
│Trans-Europe Express             │Owen Hatherley  │
╰─────────────────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_withSort(t *testing.T) {
	m := booksModel()
	m.sort = func(a, b book) int {
		// sort by author then by their book
		cmp := strings.Compare(a.author, b.author)
		if cmp == 0 {
			return strings.Compare(a.title, b.title)
		}
		return cmp
	}
	want := strings.TrimSpace(`
╭────────────────┬─────────────────────────────────╮
│Emile Zola      │Nana                             │
│Ernest Hemingway│A Farewell to Arms               │
│Ernest Hemingway│The Sun Also Rises               │
│Ernest Hemingway│To Have and Have Not             │
│Martin Amis     │London Fields                    │
│Martin Amis     │Money                            │
│Marx & Engels   │The German Ideology              │
│Owen Hatherley  │Trans-Europe Express             │
│Thomas Mann     │Death in Venice and Other Stories│
│Ulysses         │James Joyce                      │
╰────────────────┴─────────────────────────────────╯
`)
	assert.Equal(t, want, m.View())
}

//func TestModel_Height(t *testing.T) {
//	m := modelWithWindowStringData()
//	m.Height(5)
//	want := strings.TrimSpace(`
//╭─────────────────────────────────┬─────────────╮
//│The German Ideology              │Marx & Engels│
//│Death in Venice and Other Stories│Thomas Mann  │
//│Money                            │Martin Amis  │
//╰─────────────────────────────────┴─────────────╯
//`)
//	assert.Equal(t, want, m.View())
//}

func TestModel_Height_0(t *testing.T) {
	m := booksModel()
	m.Height(0)
	want := strings.TrimSpace(`
╭───────────────────┬─────────────╮
│The German Ideology│Marx & Engels│
╰───────────────────┴─────────────╯
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

func booksModel() Model[book] {
	m := New[book]()
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
