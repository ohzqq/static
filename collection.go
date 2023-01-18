package static

import (
	"github.com/ohzqq/fidi"
)

type Collection struct {
	*Page
}

func NewCollection(path string) Collection {
	tree := fidi.NewTree(path)
	col := Collection{
		Page: NewPage(tree),
	}
	col.root = path
	col.GetChildren()

	for _, page := range col.Children {
		page.GetChildren()
	}

	return col
}
