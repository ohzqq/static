package static

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ohzqq/fidi"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

type Page struct {
	fidi.Tree
	Title        string
	Css          []string
	Scripts      []string
	Color        map[string]string
	Html         Html
	buildOpts    []BuildOpt
	HtmlFiles    []fidi.File
	hasIndex     bool
	FullNav      bool
	regen        bool
	gen          bool
	index        fidi.File
	Files        []fidi.File
	Assets       []Asset
	Children     []*Page
	isCollection bool
	Nav          []map[string]any
	Breadcrumbs  []map[string]any
	tmpl         *template.Template
	root         string
	profile      string
	Opts         *Builder
}

type BuildOpt func(p *Page)

func NewPage(dir fidi.Tree, opts ...BuildOpt) *Page {
	page := Page{
		Tree:      dir,
		Css:       GetCss("global"),
		Scripts:   GetScripts("global"),
		Html:      GetHtml("global"),
		profile:   "global",
		Color:     viper.GetStringMapString("color"),
		Opts:      &Builder{},
		buildOpts: opts,
	}
	page.HtmlFiles = page.FilterByExt(".html")
	page.Index()
	page.Files = page.Leaves()

	if dir.Info().Rel() == "." {
		page.Title = "Home"
	} else {
		page.Title = dir.Info().Base
	}

	return &page
}

func (page *Page) BuildOpts(opts ...BuildOpt) *Page {
	page.buildOpts = append(page.buildOpts, opts...)
	return page
}

func (p *Page) Build() {
	fmt.Printf("building %s\n", p.Info().Name)
	for _, opt := range p.buildOpts {
		opt(p)
	}
	p.Render()
}

func Input(tree fidi.Tree) BuildOpt {
	return func(p *Page) {
		p.Tree = tree
		p.HtmlFiles = p.FilterByExt(".html")
		p.Files = p.Leaves()
		p.Index()
		if tree.Info().Rel() == "." {
			p.Title = "Home"
		} else {
			p.Title = tree.Info().Base
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
		//p.FullNav = full
		p.Breadcrumbs = getBreadcrumbs(p.Tree)
		p.Nav = getNav(p)
	}
}

func Breadcrumbs(tree fidi.Tree) BuildOpt {
	return func(p *Page) {
	}
}

func Collection() BuildOpt {
	return func(page *Page) {
		for _, dir := range page.Tree.Children() {
			p := NewPage(dir, page.buildOpts...)
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

func (p *Page) Index() *Page {
	for _, file := range p.HtmlFiles {
		if file.Base == "index.html" {
			p.index = file
			p.hasIndex = true
		}
	}
	return p
}

func (p Page) HasIndex() bool {
	for _, file := range p.HtmlFiles {
		if file.Base == "index.html" {
			return true
		}
	}
	return false
}

func (p *Page) SetTmpl(tmpl *template.Template) *Page {
	p.tmpl = tmpl
	return p
}

func (p Page) Render() string {
	if p.gen {
		tmpl := Templates.Lookup("base")
		name := filepath.Join(p.Info().Path(), "index.html")

		file, err := os.Create(name)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		err = tmpl.Execute(file, p)
		if err != nil {
			log.Fatal(err)
		}
		return name
	}

	return ""
}

func (p Page) Content() string {
	var buf bytes.Buffer
	err := p.tmpl.Execute(&buf, p)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (page *Page) GetChildren() []*Page {
	for _, dir := range page.Tree.Children() {
		p := NewPage(dir, page.buildOpts...)
		page.Children = append(page.Children, p)
	}
	return page.Children
}

func (p *Page) NewAsset(file fidi.File) *Page {
	asset := Asset{
		File:       file,
		Html:       p.Html,
		Attributes: make(map[string]any),
	}
	p.Assets = append(p.Assets, asset)
	return p
}

func (p Page) Url() map[string]any {
	url := make(map[string]any)
	url["depth"] = p.Info().Depth
	url["href"] = p.RelUrl()
	url["text"] = p.Title
	return url
}

func (p Page) AbsUrl() string {
	url := "/"
	url += filepath.Join(p.Info().Rel(), "index.html")
	return url
}

func (p Page) RelUrl() string {
	return "." + p.AbsUrl()
}

func (p *Page) FilterByExt(ext ...string) []fidi.File {
	return p.Filter(fidi.ExtFilter(ext...))
}

func (p *Page) FilterByMime(mime ...string) []fidi.File {
	return p.Filter(fidi.MimeFilter(mime...))
}

func (p Page) ReadCss() []string {
	return p.Css
}

func (p Page) ReadScripts() []string {
	return p.Scripts
}
