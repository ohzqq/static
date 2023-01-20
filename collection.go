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
	fmt.Println(p.FullNav)
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

func getFiles(page *Page, rel string) []map[string]any {
	var files []map[string]any
	for _, file := range page.Files {
		url := map[string]any{
			"href":  filepath.Join(rel, file.Base),
			"text":  file.Base,
			"depth": file.Depth,
		}

		files = append(files, url)
	}
	return files
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
		if p.FullNav {
			url["files"] = getFiles(p, rel)
		}

		nav = append(nav, url)
	}

	for idx, d := range lo.Uniq(depth) {
		for _, n := range nav {
			if n["depth"].(int) == d {
				n["depth"] = idx
			}
			if page.FullNav {
				files := n["files"].([]map[string]any)
				if len(files) > 0 {
					for _, f := range files {
						if f["depth"].(int) == d {
							f["depth"] = idx
						}
					}
				}
			}
		}
	}

	return nav
}
