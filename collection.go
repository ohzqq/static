package static

import (
	"github.com/ohzqq/fidi"
)

type Collection struct {
	Page
}

func NewCollection(path string) Collection {
	tree := fidi.NewTree(path)
	col := Collection{
		Page: NewPage(tree),
	}
	col.GetIndex()

	return col
}

func (c Collection) Pages() []Page {
	var pages []Page
	for _, dir := range c.Children() {
		page := NewPage(dir)
		page.GetIndex()
		if page.HasIndex {
			pages = append(pages, page)
		}
	}
	return pages
}
