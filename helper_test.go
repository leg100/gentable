package gentable

type isbn string

type book struct {
	isbn   isbn
	title  string
	author string
}

func renderBook(b book) []string {
	return []string{b.author, b.title}
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

//func main() {
//	t := table.New().Headers("title", "author").
//		Width(30).
//		Height(8).
//		DisableOverflowRow()
//	for _, book := range books {
//		t.Row(book.title, book.author)
//	}
//	t.Offset(2)
//	fmt.Println(t.String())
//}
