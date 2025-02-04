package gentable

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	m := modelWithStringData()
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
	data := newData(
		func(v []string) ID {
			id := i
			i++
			return id
		},
		func(v []string) []string { return v },
	)
	data.sort = func(a, b []string) int {
		// sort by first column of a and b, then by second column, etc
		cols := min(len(a), len(b))
		for i := range cols {
			switch {
			case a[i] < b[i]:
				return -1
			case a[i] > b[i]:
				return 1
			default:
				continue
			}
		}
		return 0
	}
	for _, book := range books {
		data.Append([]string{
			book.author,
			book.title,
		})
	}
	m := NewDefault(data)
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
	m := modelWithStringData()
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

func TestModel_WithGenData(t *testing.T) {
	data := newData(
		func(b book) ID {
			return b.title
		},
		func(b book) []string {
			return []string{
				b.title,
				b.author,
			}
		},
	)
	m := NewDefault(data)
	for _, bk := range books {
		data.Append(book{
			title:  bk.title,
			author: bk.author,
		})
	}
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

var i = 0

func modelWithStringData() Model {
	data := newData(
		func(v []string) ID {
			id := i
			i++
			return id
		},
		func(v []string) []string { return v },
	)
	for _, book := range books {
		data.Append([]string{
			book.title,
			book.author,
		})
	}
	return NewDefault(data)
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
