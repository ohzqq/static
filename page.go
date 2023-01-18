package static

import (
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
	Nav      []map[string]any
	Items    []fidi.File
	Children []*Page
	Url      string
	root     string
}

func NewPage(dir fidi.Tree) *Page {
	page := Page{
		Tree:    dir,
		Css:     GetCss("global"),
		Scripts: GetScripts("global"),
		Html:    GetHtml("global"),
		Color:   viper.GetStringMapString("color"),
	}

	url := page.url()
	if dir.Info().Rel() == "." {
		page.Title = "Home"
		url["depth"] = 1
		url["text"] = page.Title
	} else {
		page.Title = dir.Info().Base
	}

	page.Nav = []map[string]any{url}

	for _, file := range page.FilterByExt(".html") {
		if file.Base == "index.html" {
			page.HasIndex = true
			page.Index = file
		}
	}

	return &page
}

func (page *Page) GetChildren() {
	for _, dir := range page.Tree.Children() {
		p := NewPage(dir)
		if p.HasIndex {
			page.Nav = append(page.Nav, p.url())
			page.Children = append(page.Children, p)
		}
	}
}

func (p *Page) Profile(pro string) *Page {
	css := GetCss(pro)
	p.Css = append(p.Css, css...)

	scripts := GetScripts(pro)
	p.Scripts = append(p.Scripts, scripts...)

	html := GetHtml(pro)
	maps.Copy(p.Html, html)

	return p
}

func (p *Page) CreateIndex() *Page {
	if !p.HasIndex {
		//Create filetree page
	}
	return p
}

func (p Page) url() map[string]any {
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

func (p Page) Breadcrumbs() []map[string]string {
	var crumbs []map[string]string

	totalP := len(p.Parents())
	for _, parent := range p.Parents() {
		totalP--

		path := ".." + strings.Repeat("/..", totalP)
		path = filepath.Join(path, "index.html")

		name := parent.Info().Base
		if parent.Info().Rel() == "." {
			name = "Home"
		}

		link := map[string]string{
			"href": path,
			"text": name,
		}
		crumbs = append(crumbs, link)
	}

	return crumbs
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
