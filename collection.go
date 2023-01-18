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
	col.root = path

	return col
}

func GetHtmlFiles(tree fidi.Tree) []fidi.File {
	return tree.Filter(fidi.ExtFilter(".html"))
}
