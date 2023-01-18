package static

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/ohzqq/fidi"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

type Page struct {
	fidi.Tree
	Title    string
	Css      []string
	Scripts  []string
	Color    map[string]string
	Html     Html
	HasIndex bool
	Index    fidi.File
	Items    []fidi.File
	Assets   []Asset
	Children []*Page
	root     string
	profile  string
}

func NewPage(dir fidi.Tree) *Page {
	page := Page{
		Tree:    dir,
		Css:     GetCss("global"),
		Scripts: GetScripts("global"),
		Html:    GetHtml("global"),
		Color:   viper.GetStringMapString("color"),
	}

	if dir.Info().Rel() == "." {
		page.Title = "Home"
	} else {
		page.Title = dir.Info().Base
	}

	for _, file := range page.FilterByExt(".html") {
		if file.Base == "index.html" {
			page.HasIndex = true
			page.Index = file
		}
	}

	return &page
}

func (p Page) Content() string {
}

func (page *Page) GetChildren() {
	for _, dir := range page.Tree.Children() {
		p := NewPage(dir)
		if p.HasIndex {
			if page.profile != "" {
				p.Profile(page.profile)
			}
			page.Children = append(page.Children, p)
		}
	}
}

func (p *Page) NewAsset(file fidi.File) *Page {
	asset := Asset{
		File: file,
		Html: p.Html,
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

func (p *Page) CreateIndex() *Page {
	if !p.HasIndex {
		//Create filetree page
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
