package static

import (
	"log"

	"github.com/ohzqq/fidi"
)

type Builder struct {
	Nav          bool
	ListAll      bool
	IndexOnly    bool
	Gen          bool
	Regen        bool
	IsCollection bool
	NoThumbs     bool
	Tmpl         string
	Profile      string
	Input        string
	Mimetypes    []string
	Exts         []string
}

func New(path string) *Builder {
	return &Builder{
		Input: path,
	}
}

func (b *Builder) Collection() *Builder {
	b.IsCollection = true
	b.Nav = true
	return b
}

func (b *Builder) Page() *Builder {
	b.IsCollection = false
	return b
}

func (b Builder) Opts() []BuildOpt {
	var opts []BuildOpt

	if b.IsCollection {
		opts = append(opts, Collection())
	}

	if b.Nav {
		opts = append(opts, Nav(b.ListAll))
	}

	switch {
	case b.Gen:
		opts = append(opts, Gen())
	case b.Regen:
		opts = append(opts, Regen())
	}

	if b.Profile != "" {
		opts = append(opts, Profile(b.Profile))
	}

	return opts
}

func (b *Builder) Build() {
	if b.Input == "" {
		log.Fatal("no input")
	}
	tree := fidi.NewTree(b.Input)

	page := NewPage(tree)
	page.Build(b.Opts()...)

	if b.IsCollection {
		for _, child := range page.Children {
			child.Build(b.Opts()...)
		}
	}
}

func Gen() BuildOpt {
	return func(pg *Page) {
		switch {
		case !pg.hasIndex:
			pg.Gen = true
		default:
			pg.Gen = false
		}
	}
}

func Regen() BuildOpt {
	return func(pg *Page) {
		switch {
		case pg.Gen:
			pg.Gen = false
		default:
			pg.Gen = true
		}
	}
}

func Nav(full bool) BuildOpt {
	return func(pg *Page) {
		pg.FullNav = full
		pg.getBreadcrumbs()
		pg.getNav()
	}
}

func Collection() BuildOpt {
	return func(pg *Page) {
		pg.setChildren()
	}
}

func Profile(pro string) BuildOpt {
	return func(pg *Page) {
		pg.SetProfile(pro)
	}
}

func AssetFilters(filters ...fidi.Filter) BuildOpt {
	return func(pg *Page) {
		pg.filters = filters
	}
}
