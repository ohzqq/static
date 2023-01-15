package static

import (
	"github.com/ohzqq/fidi"
)

type Collection struct {
	fidi.Tree
}

func NewCollection(path string) Collection {
	return Collection{
		Tree: fidi.NewTree(path),
	}
}

func (c Collection) Pages() []Page {
	var pages []Page
	for _, dir := range c.Nodes {
		page := NewPage(dir)
		page.FilterByExt(".html")
		pages = append(pages, page)
	}
	return pages
}

func GetNodesByDepth(tree fidi.Tree, d int) []fidi.Dir {
	var nodes []fidi.Dir
	for _, node := range tree.Nodes {
		if node.Depth == d {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func GetChildrenByDepth(tree fidi.Tree, d int) []fidi.Dir {
	if d == 0 {
		return tree.Nodes
	}

	var nodes []fidi.Dir
	for i := d + 2; i < len(tree.Nodes); i++ {
		n := GetNodesByDepth(tree, i)
		nodes = append(nodes, n...)
	}
	return nodes
}

func GetParentsByDepth(tree fidi.Tree, d int) []fidi.Dir {
	if d == 0 {
		return tree.Nodes
	}

	var nodes []fidi.Dir
	for i := d; i > 0; i-- {
		n := GetNodesByDepth(tree, i)
		nodes = append(nodes, n...)
	}
	return nodes
}
