package static

import (
	"fmt"

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
	fmt.Printf("opts %s\n", len(opts))
	col.root = path

	col.GetChildren()

	for _, page := range col.Children {
		page.GetChildren()
	}

	return col
}
