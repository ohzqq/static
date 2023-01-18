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
	Nav      []map[string]any
	Items    []fidi.File
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

	url := page.Url()
	if dir.Info().Rel() == "." {
		page.Title = "Home"
		url["depth"] = 1
	} else {
		page.Title = dir.Info().Base
	}
	url["text"] = page.Title

	//page.Nav = []map[string]any{url}

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
			if page.profile != "" {
				p.Profile(page.profile)
			}
			rel, err := filepath.Rel(page.Info().Rel(), p.Info().Rel())
			if err != nil {
				log.Fatal(err)
			}
			url := map[string]any{
				"href":  filepath.Join(rel, "index.html"),
				"text":  p.Title,
				"depth": p.Info().Depth,
			}
			page.Nav = append(page.Nav, url)
			page.Children = append(page.Children, p)
		}
	}
}

func (p *Page) Profile(pro string) *Page {
	p.profile = pro

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
