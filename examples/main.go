package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/leg100/gentable"
)

func main() {
	m := model{
		Model: gentable.New[book](),
	}
	m.Headers("isbn", "title", "author")
	m.Height(6)
	m.Width(30)
	for _, bk := range books {
		m.Append(gentable.Row[book]{
			V: bk,
			Cells: []string{
				string(bk.isbn),
				bk.title,
				bk.author,
			},
		})
	}

	p := tea.NewProgram(m)
	_, _ = p.Run()

}

type model struct {
	gentable.Model[book]
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	body := strings.Builder{}

	fmt.Fprintf(&body, "cursor: %d\n", m.Model.CursorIndex())
	fmt.Fprintf(&body, "start: %d\n", m.Model.StartIndex())
	fmt.Fprintf(&body, "size: %d\n", m.Model.Size())

	body.WriteString(m.Model.View())

	return body.String()
}

type isbn string

type book struct {
	isbn   isbn
	title  string
	author string
}

var books = []book{
	{
		title:  "The German Ideology",
		author: "Marx & Engels",
		isbn:   "390039233003",
	},
	{
		title:  "Death in Venice and Other Stories",
		author: "Thomas Mann",
		isbn:   "390039233004",
	},
	{
		title:  "Money",
		author: "Martin Amis",
		isbn:   "390039233005",
	},
	{
		title:  "London Fields",
		author: "Martin Amis",
		isbn:   "390039233006",
	},
	{
		title:  "Nana",
		author: "Emile Zola",
		isbn:   "390039233007",
	},
	{
		title:  "To Have and Have Not",
		author: "Ernest Hemingway",
		isbn:   "390039233008",
	},
	{
		title:  "The Sun Also Rises",
		author: "Ernest Hemingway",
		isbn:   "390039233009",
	},
	{
		title:  "A Farewell to Arms",
		author: "Ernest Hemingway",
		isbn:   "390039233010",
	},
	{
		title:  "James Joyce",
		author: "Ulysses",
		isbn:   "390039233011",
	},
	{
		title:  "Trans-Europe Express",
		author: "Owen Hatherley",
		isbn:   "390039233012",
	},
}
