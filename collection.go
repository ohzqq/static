package static

import (
	"github.com/ohzqq/fidi"
)

type Collection struct {
	*Page
	Pages []*Page
}

func NewCollection(path string) Collection {
	tree := fidi.NewTree(path)
	col := Collection{
		Page: NewPage(tree),
	}
	col.Title = "Home"

	col.Url = map[string]any{
		"href":  "./index.html",
		"text":  "Home",
		"depth": 1,
	}
	col.Nav = append(col.Nav, col.Url)

	for _, dir := range col.Children() {
		page := NewPage(dir)
		if page.HasIndex {
			page.Url["depth"] = page.Info().Depth
			col.Nav = append(col.Nav, page.Url)
			col.Pages = append(col.Pages, page)
		}
	}

	for _, page := range col.Pages {
		page.Nav = col.Nav
	}

	return col
}
