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
	pages    []fidi.File
	Url      map[string]any
	Nav      []map[string]any
	Items    []fidi.File
	root     string
}

func NewPage(dir fidi.Tree) *Page {
	page := Page{
		Tree:    dir,
		Css:     GetCss("global"),
		Scripts: GetScripts("global"),
		Html:    GetHtml("global"),
		Color:   viper.GetStringMapString("color"),
		Title:   dir.Info().Base,
		Url:     make(map[string]any),
	}

	files := page.Filter(fidi.ExtFilter(".html"))
	for _, file := range files {
		if file.Base == "index.html" {
			page.HasIndex = true
			page.Index = file
			page.Url["href"] = "./" + file.Rel()
			page.Url["text"] = page.Title
			page.pages = append(page.pages, file)
		}
	}

	return &page
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
