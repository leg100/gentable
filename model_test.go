package gentable

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	m := booksModel()
	want := strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233003 The German Ideology               Marx & Engels   
390039233004 Death in Venice and Other Stories Thomas Mann     
390039233005 Money                             Martin Amis     
390039233006 London Fields                     Martin Amis     
390039233007 Nana                              Emile Zola      
390039233008 To Have and Have Not              Ernest Hemingway
390039233009 The Sun Also Rises                Ernest Hemingway
390039233010 A Farewell to Arms                Ernest Hemingway
390039233011 James Joyce                       Ulysses         
390039233012 Trans-Europe Express              Owen Hatherley  
`, "\n")
	assert.Equal(t, want, m.View())
}

func TestModel_Sort(t *testing.T) {
	m := booksModel()
	m.SetSortFunc(func(a, b book) int {
		// sort by author then by their book
		cmp := strings.Compare(a.author, b.author)
		if cmp == 0 {
			return strings.Compare(a.title, b.title)
		}
		return cmp
	})
	want := strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233003 The German Ideology               Marx & Engels   
390039233004 Death in Venice and Other Stories Thomas Mann     
390039233005 Money                             Martin Amis     
390039233006 London Fields                     Martin Amis     
390039233007 Nana                              Emile Zola      
390039233008 To Have and Have Not              Ernest Hemingway
390039233009 The Sun Also Rises                Ernest Hemingway
390039233010 A Farewell to Arms                Ernest Hemingway
390039233011 James Joyce                       Ulysses         
390039233012 Trans-Europe Express              Owen Hatherley  
`, "\n")
	assert.Equal(t, want, m.View())
}

func TestModel_Filter(t *testing.T) {
	m := booksModel()
	m.ApplyFilter(func(b book) bool {
		return b.author == "Ernest Hemingway"
	})
	want := strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233008 To Have and Have Not              Ernest Hemingway
390039233009 The Sun Also Rises                Ernest Hemingway
390039233010 A Farewell to Arms                Ernest Hemingway
`, "\n")
	assert.Equal(t, want, m.View())

	m.RemoveFilter()

	want = strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233003 The German Ideology               Marx & Engels   
390039233004 Death in Venice and Other Stories Thomas Mann     
390039233005 Money                             Martin Amis     
390039233006 London Fields                     Martin Amis     
390039233007 Nana                              Emile Zola      
390039233008 To Have and Have Not              Ernest Hemingway
390039233009 The Sun Also Rises                Ernest Hemingway
390039233010 A Farewell to Arms                Ernest Hemingway
390039233011 James Joyce                       Ulysses         
390039233012 Trans-Europe Express              Owen Hatherley  
`, "\n")
	assert.Equal(t, want, m.View())
}

func TestModel_Height_Before(t *testing.T) {
	m := booksModel()
	m.Height(5)
	want := strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233003 The German Ideology               Marx & Engels   
390039233004 Death in Venice and Other Stories Thomas Mann     
390039233005 Money                             Martin Amis     
390039233006 London Fields                     Martin Amis     
`, "\n")
	assert.Equal(t, want, m.View())
}

func TestModel_Height_After(t *testing.T) {
	m := booksModel()
	m.Height(5)
	want := strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233003 The German Ideology               Marx & Engels   
390039233004 Death in Venice and Other Stories Thomas Mann     
390039233005 Money                             Martin Amis     
390039233006 London Fields                     Martin Amis     
`, "\n")
	assert.Equal(t, want, m.View())
}

func TestModel_Height_0(t *testing.T) {
	m := booksModel()
	m.Height(0)
	want := strings.Trim(`
ISBN         TITLE                             AUTHOR          
`, "\n")
	assert.Equal(t, want, m.View())
}

func TestModel_PageDown(t *testing.T) {
	m := booksModel()
	m.Height(5)

	m.window.PageDown()

	want := strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233004 Death in Venice and Other Stories Thomas Mann     
390039233005 Money                             Martin Amis     
390039233006 London Fields                     Martin Amis     
390039233007 Nana                              Emile Zola      
`, "\n")
	assert.Equal(t, want, m.View())

	m.window.PageDown()

	want = strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233008 To Have and Have Not              Ernest Hemingway
390039233009 The Sun Also Rises                Ernest Hemingway
390039233010 A Farewell to Arms                Ernest Hemingway
390039233011 James Joyce                       Ulysses         
`, "\n")
	assert.Equal(t, want, m.View())

	m.window.PageDown()

	want = strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233009 The Sun Also Rises                Ernest Hemingway
390039233010 A Farewell to Arms                Ernest Hemingway
390039233011 James Joyce                       Ulysses         
390039233012 Trans-Europe Express              Owen Hatherley  
`, "\n")
	assert.Equal(t, want, m.View())

	// no more books so should be a no-op
	m.window.PageDown()

	want = strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233009 The Sun Also Rises                Ernest Hemingway
390039233010 A Farewell to Arms                Ernest Hemingway
390039233011 James Joyce                       Ulysses         
390039233012 Trans-Europe Express              Owen Hatherley  
`, "\n")
	assert.Equal(t, want, m.View())
}

func TestModel_PageUp(t *testing.T) {
	m := booksModel()
	m.Height(5)
	m.window.PageDown()
	m.window.PageDown()
	m.window.PageUp()
	want := strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233007 Nana                              Emile Zola      
390039233008 To Have and Have Not              Ernest Hemingway
390039233009 The Sun Also Rises                Ernest Hemingway
390039233010 A Farewell to Arms                Ernest Hemingway
`, "\n")
	assert.Equal(t, want, m.View())
}

func TestModel_BottomTop(t *testing.T) {
	m := booksModel()
	m.Height(5)

	m.window.toBottom()

	want := strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233009 The Sun Also Rises                Ernest Hemingway
390039233010 A Farewell to Arms                Ernest Hemingway
390039233011 James Joyce                       Ulysses         
390039233012 Trans-Europe Express              Owen Hatherley  
`, "\n")
	assert.Equal(t, want, m.View())

	m.window.toTop()

	want = strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233003 The German Ideology               Marx & Engels   
390039233004 Death in Venice and Other Stories Thomas Mann     
390039233005 Money                             Martin Amis     
390039233006 London Fields                     Martin Amis     
`, "\n")
	assert.Equal(t, want, m.View())
}

func TestModel_Expand(t *testing.T) {
	m := booksModel()
	m.Height(5)

	want := strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233003 The German Ideology               Marx & Engels   
390039233004 Death in Venice and Other Stories Thomas Mann     
390039233005 Money                             Martin Amis     
390039233006 London Fields                     Martin Amis     
`, "\n")
	assert.Equal(t, want, m.View())

	m.Height(8)

	want = strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233003 The German Ideology               Marx & Engels   
390039233004 Death in Venice and Other Stories Thomas Mann     
390039233005 Money                             Martin Amis     
390039233006 London Fields                     Martin Amis     
390039233007 Nana                              Emile Zola      
390039233008 To Have and Have Not              Ernest Hemingway
390039233009 The Sun Also Rises                Ernest Hemingway
`, "\n")
	assert.Equal(t, want, m.View())
}

func TestModel_Shrink(t *testing.T) {
	m := booksModel()
	m.Height(8)

	want := strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233003 The German Ideology               Marx & Engels   
390039233004 Death in Venice and Other Stories Thomas Mann     
390039233005 Money                             Martin Amis     
390039233006 London Fields                     Martin Amis     
390039233007 Nana                              Emile Zola      
390039233008 To Have and Have Not              Ernest Hemingway
390039233009 The Sun Also Rises                Ernest Hemingway
`, "\n")
	assert.Equal(t, want, m.View())

	m.Height(5)

	want = strings.Trim(`
ISBN         TITLE                             AUTHOR          
390039233003 The German Ideology               Marx & Engels   
390039233004 Death in Venice and Other Stories Thomas Mann     
390039233005 Money                             Martin Amis     
390039233006 London Fields                     Martin Amis     
`, "\n")
	assert.Equal(t, want, m.View())

}

func booksModel() Model[book] {
	m := New[book]("ISBN", "TITLE", "AUTHOR")
	m.Height(30)
	m.Width(63)
	m.disableCursorHighlighting = true
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
