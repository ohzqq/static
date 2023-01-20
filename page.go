package static

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ohzqq/fidi"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

type Page struct {
	fidi.Tree
	Title     string
	Css       []string
	Scripts   []string
	Color     map[string]string
	Html      Html
	buildOpts []BuildOpt
	HtmlFiles []fidi.File
	hasIndex  bool
	regen     bool
	gen       bool
	index     fidi.File
	Items     []fidi.File
	Assets    []Asset
	Children  []*Page
	tmpl      *template.Template
	root      string
	profile   string
}

type BuildOpt func(p *Page) BuildOpt

func NewPage(dir fidi.Tree, opts ...BuildOpt) *Page {
	page := Page{
		Tree:      dir,
		Css:       GetCss("global"),
		Scripts:   GetScripts("global"),
		Html:      GetHtml("global"),
		profile:   "global",
		Color:     viper.GetStringMapString("color"),
		buildOpts: opts,
	}
	page.HtmlFiles = page.FilterByExt(".html")
	page.Index()

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

func Gen() BuildOpt {
	return func(p *Page) BuildOpt {
		switch {
		case !p.hasIndex:
			p.gen = true
		default:
			p.gen = false
		}
		return Gen()
	}
}

func Regen() BuildOpt {
	return func(p *Page) BuildOpt {
		switch {
		case p.gen:
			p.gen = false
		default:
			p.gen = true
		}
		return Regen()
	}
}

func Profile(pro string) BuildOpt {
	return func(p *Page) BuildOpt {
		p.tmpl = GetTemplate(pro)
		p.profile = pro
		if pro == "global" {
			p.Css = GetCss("global")
			p.Scripts = GetScripts("global")
			p.Html = GetHtml("global")
			//p.profile = "global"
		} else {
			css := GetCss(pro)
			p.Css = append(p.Css, css...)

			scripts := GetScripts(pro)
			p.Scripts = append(p.Scripts, scripts...)

			html := GetHtml(pro)
			maps.Copy(p.Html, html)
		}

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

		return Profile("global")
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

func (page *Page) GetChildren() {
	for _, dir := range page.Tree.Children() {
		p := NewPage(dir, page.buildOpts...)
		page.Children = append(page.Children, p)
	}
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

func (p *Page) Profile(pro string) *Page {
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
	p.Items = items

	for _, i := range items {
		p.NewAsset(i)
	}

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

func (p Page) Breadcrumbs() []map[string]any {
	var crumbs []map[string]any

	totalP := len(p.Parents())
	for _, parent := range p.Parents() {
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

func (page *Page) Nav() []map[string]any {
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

func (p *Page) FilterByExt(ext ...string) []fidi.File {
	//p.Items = p.Filter(fidi.ExtFilter(ext...))
	//return p
	return p.Filter(fidi.ExtFilter(ext...))
}

func (p *Page) FilterByMime(mime ...string) []fidi.File {
	return p.Filter(fidi.MimeFilter(mime...))
	//p.Items = p.Filter(fidi.MimeFilter(mime...))
	//return p
}

func (p Page) ReadCss() []string {
	return p.Css
}

func (p Page) ReadScripts() []string {
	return p.Scripts
}
