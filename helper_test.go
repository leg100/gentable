package gentable

var books = []struct {
	title  string
	author string
}{
	{
		title:  "The German Ideology",
		author: "Marx & Engels",
	},
	{
		title:  "Death in Venice and Other Stories",
		author: "Thomas Mann",
	},
	{
		title:  "Money",
		author: "Martin Amis",
	},
	{
		title:  "London Fields",
		author: "Martin Amis",
	},
	{
		title:  "Nana",
		author: "Emile Zola",
	},
	{
		title:  "To Have and Have Not",
		author: "Ernest Hemingway",
	},
	{
		title:  "The Sun Also Rises",
		author: "Ernest Hemingway",
	},
	{
		title:  "A Farewell to Arms",
		author: "Ernest Hemingway",
	},
	{
		title:  "James Joyce",
		author: "Ulysses",
	},
	{
		title:  "Trans-Europe Express",
		author: "Owen Hatherley",
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
