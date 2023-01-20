package static

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/ohzqq/fidi"
	"github.com/samber/lo"
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
	cfgCollectionPage(col.Page)

	return col
}

func (col Collection) Build() {
	col.Page.Build()
	for _, page := range col.Children {
		page.Build()
	}
}

func cfgCollectionPage(p *Page) *Page {
	p.GetChildren()
	p.Breadcrumbs = getBreadcrumbs(p.Tree)
	p.Nav = getNav(p)

	for _, page := range p.Children {
		cfgCollectionPage(page)
	}

	return p
}

func (col *Collection) getChildren() *Collection {
	col.GetChildren()
	for _, page := range col.Children {
		page.GetChildren()
	}
	return col
}

func getBreadcrumbs(tree fidi.Tree) []map[string]any {
	var crumbs []map[string]any

	totalP := len(tree.Parents())
	for _, parent := range tree.Parents() {
		totalP--

		path := ".." + strings.Repeat("/..", totalP)
		path = filepath.Join(path, "index.html")

		name := parent.Info().Base
		if parent.Info().Rel() == "." {
			name = "Home"
		}

		link := map[string]any{
			"href":  path,
			"text":  name,
			"depth": parent.Info().Depth,
		}
		crumbs = append(crumbs, link)
	}

	return crumbs
}

func getNav(page *Page) []map[string]any {
	var depth []int
	var nav []map[string]any
	for _, p := range page.Children {
		self := page.Info().Rel()
		child := p.Info().Rel()
		rel, err := filepath.Rel(self, child)
		if err != nil {
			log.Fatal(err)
		}

		depth = append(depth, p.Info().Depth)

		url := map[string]any{
			"href":  filepath.Join(rel, "index.html"),
			"text":  p.Title,
			"depth": p.Info().Depth,
		}

		nav = append(nav, url)
	}

	depth = lo.Uniq(depth)
	fmt.Printf("depht %v\n", depth)

	return nav
}
