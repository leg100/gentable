package gentable

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	m := New()
	for _, book := range books {
		m.Rows([]string{
			book.title,
			book.author,
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

func TestModel_Height(t *testing.T) {
	m := New()
	m.Height(5)
	for _, book := range books {
		m.Rows([]string{
			book.title,
			book.author,
		})
	}
	want := strings.TrimSpace(`
╭─────────────────────────────────┬─────────────╮
│The German Ideology              │Marx & Engels│
│Death in Venice and Other Stories│Thomas Mann  │
│Money                            │Martin Amis  │
╰─────────────────────────────────┴─────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_Height_0(t *testing.T) {
	m := New()
	m.Height(0)
	for _, book := range books {
		m.Rows([]string{
			book.title,
			book.author,
		})
	}
	want := strings.TrimSpace(`
╭───────────────────┬─────────────╮
│The German Ideology│Marx & Engels│
╰───────────────────┴─────────────╯
`)
	assert.Equal(t, want, m.View())
}

func TestModel_PageUp(t *testing.T) {
	m := New()
	m.Height(5)
	for _, book := range books {
		m.Rows([]string{
			book.title,
			book.author,
		})
	}

	m.data.PageDown()
	want := strings.TrimSpace(`
╭────────────────────┬────────────────╮
│London Fields       │Martin Amis     │
│Nana                │Emile Zola      │
│To Have and Have Not│Ernest Hemingway│
╰────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

	m.data.PageDown()
	want = strings.TrimSpace(`
╭──────────────────┬────────────────╮
│The Sun Also Rises│Ernest Hemingway│
│A Farewell to Arms│Ernest Hemingway│
│James Joyce       │Ulysses         │
╰──────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

	m.data.PageDown()
	want = strings.TrimSpace(`
╭────────────────────┬──────────────╮
│Trans-Europe Express│Owen Hatherley│
│                    │              │
│                    │              │
╰────────────────────┴──────────────╯
`)
	assert.Equal(t, want, m.View())

	m.data.PageDown()
	want = strings.TrimSpace(`
╭────────────────────┬──────────────╮
│Trans-Europe Express│Owen Hatherley│
│                    │              │
│                    │              │
╰────────────────────┴──────────────╯
`)
	assert.Equal(t, want, m.View())

	m.data.PageUp()
	want = strings.TrimSpace(`
╭──────────────────┬────────────────╮
│The Sun Also Rises│Ernest Hemingway│
│A Farewell to Arms│Ernest Hemingway│
│James Joyce       │Ulysses         │
╰──────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

	m.data.PageUp()
	want = strings.TrimSpace(`
╭────────────────────┬────────────────╮
│London Fields       │Martin Amis     │
│Nana                │Emile Zola      │
│To Have and Have Not│Ernest Hemingway│
╰────────────────────┴────────────────╯
`)
	assert.Equal(t, want, m.View())

	m.data.PageUp()
	want = strings.TrimSpace(`
╭─────────────────────────────────┬─────────────╮
│The German Ideology              │Marx & Engels│
│Death in Venice and Other Stories│Thomas Mann  │
│Money                            │Martin Amis  │
╰─────────────────────────────────┴─────────────╯
`)
	assert.Equal(t, want, m.View())

	m.data.PageUp()
	want = strings.TrimSpace(`
╭─────────────────────────────────┬─────────────╮
│The German Ideology              │Marx & Engels│
│Death in Venice and Other Stories│Thomas Mann  │
│Money                            │Martin Amis  │
╰─────────────────────────────────┴─────────────╯
`)
	assert.Equal(t, want, m.View())
}
