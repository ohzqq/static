package static

import (
	"github.com/ohzqq/fidi"
)

type Collection struct {
	*Page
}

func NewCollection(path string, opts ...BuildOpt) Collection {
	tree := fidi.NewTree(path)

	col := Collection{
		Page: NewPage(tree),
	}
	col.BuildOpts(opts...)
	col.root = path

	col.GetChildren()

	for _, page := range col.Children {
		page.GetChildren()
	}

	return col
}

func (col Collection) Build() {
	col.Page.Build()
	for _, page := range col.Children {
		page.Build()
	}
}
