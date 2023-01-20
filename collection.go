package static

import (
	"log"
	"path/filepath"
	"strings"

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
	var nav []map[string]any

	for _, p := range page.Children {
		self := page.Info().Rel()
		child := p.Info().Rel()
		rel, err := filepath.Rel(self, child)
		if err != nil {
			log.Fatal(err)
		}

		url := map[string]any{
			"href":  filepath.Join(rel, "index.html"),
			"text":  p.Title,
			"depth": p.Info().Depth,
		}

		nav = append(nav, url)
	}

	return nav
}

func (col *Collection) getBreadcrumbs() *Collection {
	totalP := len(col.Parents())
	for _, parent := range col.Parents() {
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
		col.Breadcrumbs = append(col.Breadcrumbs, link)
	}

	return col
}

func (col *Collection) getNav() *Collection {
	for _, p := range col.Children {
		self := col.Info().Rel()
		child := p.Info().Rel()
		rel, err := filepath.Rel(self, child)
		if err != nil {
			log.Fatal(err)
		}
		url := map[string]any{
			"href":  filepath.Join(rel, "index.html"),
			"text":  p.Title,
			"depth": p.Info().Depth,
		}
		col.Nav = append(col.Nav, url)
	}
	return col
}
