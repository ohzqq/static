package static

import (
	"log"

	"github.com/ohzqq/fidi"
	"github.com/spf13/viper"
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

func (b Builder) State() string {
	var state string
	switch {
	case b.IndexOnly:
		switch {
		}
	}
	return state
}

func listAll() bool {
	switch {
	case viper.GetBool("build.all"):
		return true
	case indexOnly():
		if recurse() {
			return false
		}
		return true
	default:
		return false
	}
}

func indexOnly() bool {
	return viper.GetBool("build.index_only")
}

func noThumbs() bool {
	return viper.GetBool("build.no_thumbs")
}

func recurse() bool {
	return viper.GetBool("build.is_collection")
}

func regen() bool {
	return viper.GetBool("build.regen")
}

func hasMimes() bool {
	return len(mimes()) > 0
}

func mimes() []string {
	return viper.GetStringSlice("build.mimes")
}

func hasExts() bool {
	return len(exts()) > 0
}

func exts() []string {
	return viper.GetStringSlice("build.exts")
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

	if b.Profile != "" {
		opts = append(opts, Profile(b.Profile))
	}

	if len(b.Mimetypes) > 0 {
		opts = append(opts, FilterAssets(fidi.MimeFilter(b.Mimetypes...)))
	}

	if len(b.Exts) > 0 {
		opts = append(opts, FilterAssets(fidi.ExtFilter(b.Exts...)))
	}

	switch {
	case b.Gen:
		opts = append(opts, Gen())
	case b.Regen:
		opts = append(opts, Regen())
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
			viper.Set("build.gen", true)
		default:
			pg.Gen = false
			viper.Set("build.gen", false)
		}
	}
}

func Regen() BuildOpt {
	return func(pg *Page) {
		switch {
		case viper.GetBool("build.gen"):
			fallthrough
		case pg.Gen:
			pg.Gen = false
			viper.Set("build.gen", false)
		default:
			pg.Gen = true
			viper.Set("build.gen", true)
		}
	}
}

func Nav(full bool) BuildOpt {
	return func(pg *Page) {
		pg.FullNav = full
		pg.setBreadcrumbs()
		pg.setNav()
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

func FilterAssets(filters ...fidi.Filter) BuildOpt {
	return func(pg *Page) {
		pg.filters = append(pg.filters, filters...)
	}
}

func NoThumbs() BuildOpt {
	return func(pg *Page) {
		pg.NoThumbs = true
	}
}
