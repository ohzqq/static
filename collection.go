package static

import (
	"github.com/ohzqq/fidi"
)

type Collection struct {
	*Page
	//Pages []*Page
}

func NewCollection(path string) Collection {
	tree := fidi.NewTree(path)
	col := Collection{
		Page: NewPage(tree),
	}
	col.Title = "Home"
	col.root = path

	return col
}

func (col *Collection) Pages() []*Page {
	var pages []*Page

	//curl := map[string]any{
	//  "href":  "./index.html",
	//  "text":  "Home",
	//  "depth": 1,
	//}
	//col.Nav = append(col.Nav, curl)

	//for _, dir := range col.Children() {
	//  page := NewPage(dir)
	//  if page.HasIndex {
	//    col.Nav = append(col.Nav, page.url())
	//    pages = append(pages, page)
	//  }
	//}

	return pages
	//return col.pages
}

func GetIndexFiles(tree fidi.Tree) []fidi.File {
	var pages []fidi.File
	files := tree.Filter(fidi.ExtFilter(".html"))
	for _, file := range files {
		if file.Base == "index.html" {
			pages = append(pages, file)
		}
	}
	return pages
}
