package static

import (
	"log"
	"text/template"

	"github.com/ohzqq/fidi"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

type Builder struct {
	Nav          bool
	FullNav      bool
	Gen          bool
	Regen        bool
	isCollection bool
	Tmpl         *template.Template
	Profile      string
	Input        string
}

func New(path string) *Builder {
	return &Builder{
		Input: path,
	}
}

func (b *Builder) Collection() *Builder {
	b.isCollection = true
	b.Nav = true
	return b
}

func (b *Builder) Page() *Builder {
	b.isCollection = false
	return b
}

func (b Builder) Opts() []BuildOpt {
	var opts []BuildOpt

	if b.isCollection {
		opts = append(opts, Collection())
	}

	if b.Nav {
		opts = append(opts, Nav(b.FullNav))
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

	if b.isCollection {
		for _, child := range page.Children {
			child.Build(b.Opts()...)
		}
	}
}

func Gen() BuildOpt {
	return func(p *Page) {
		switch {
		case !p.hasIndex:
			p.gen = true
		default:
			p.gen = false
		}
	}
}

func Regen() BuildOpt {
	return func(p *Page) {
		switch {
		case p.gen:
			p.gen = false
		default:
			p.gen = true
		}
	}
}

func Nav(full bool) BuildOpt {
	return func(p *Page) {
		p.FullNav = full
		p.getBreadcrumbs()
		p.getNav()
	}
}

func Collection() BuildOpt {
	return func(page *Page) {
		for _, dir := range page.Tree.Children() {
			p := NewPage(dir)
			page.Children = append(page.Children, p)
		}
	}
}

func Profile(pro string) BuildOpt {
	return func(p *Page) {
		p.tmpl = GetTemplate(pro)
		p.profile = pro

		css := GetCss(pro)
		p.Css = append(p.Css, css...)

		scripts := GetScripts(pro)
		p.Scripts = append(p.Scripts, scripts...)

		html := GetHtml(pro)
		maps.Copy(p.Html, html)

		mt := pro + ".mime"
		ext := pro + ".ext"
		var items []fidi.File
		switch {
		case viper.IsSet(mt):
			mimes := viper.GetStringSlice(mt)
			items = p.FilterByMime(mimes...)
		case viper.IsSet(ext):
			exts := viper.GetStringSlice(ext)
			items = p.FilterByExt(exts...)
		}

		for _, i := range items {
			p.NewAsset(i)
		}
	}
}
