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
	//Profile
	Title    string
	Css      []string
	Scripts  []string
	Color    map[string]string
	Html     Html
	HasIndex bool
	Index    fidi.File
	Nav      []map[string]any
	Items    []fidi.File
	pages    []fidi.File
	Url      string
	root     string
}

func newPage(file fidi.File) *Page {
	page := Page{
		Index:   file,
		Css:     GetCss("global"),
		Scripts: GetScripts("global"),
		Html:    GetHtml("global"),
		Color:   viper.GetStringMapString("color"),
		Title:   file.Base,
		Url:     file.Rel(),
	}
	return &page
}

func NewPage(dir fidi.Tree) *Page {
	//page := newPage()
	var page *Page

	//files := page.Filter(fidi.ExtFilter(".html"))
	files := GetIndexFiles(dir)
	for _, file := range files {
		if i := file.Rel(); i == "index.html" {
			page = newPage(file)
			page.Tree = dir
		} else {
			page.pages = append(page.pages, file)
		}
	}

	return page
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

func (p *Page) Pages() []*Page {
	var pages []*Page

	curl := map[string]any{
		"href":  "./index.html",
		"text":  p.Title,
		"depth": 0,
	}
	p.Nav = append(p.Nav, curl)

	for _, dir := range p.Children() {
		page := NewPage(dir)
		if page.HasIndex {
			p.Nav = append(p.Nav, page.url())
			pages = append(pages, page)
		}
	}

	return pages
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

func (p *Page) FilterByExt(ext ...string) *Page {
	p.Items = p.Filter(fidi.ExtFilter(ext...))
	return p
}

func (p *Page) FilterByMime(mime ...string) *Page {
	p.Items = p.Filter(fidi.MimeFilter(mime...))
	return p
}

func (p Page) ReadCss() []string {
	return p.Css
}

func (p Page) ReadScripts() []string {
	return p.Scripts
}
