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
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

type Page struct {
	fidi.Tree
	Title       string
	Css         []string
	Scripts     []string
	Color       map[string]string
	Html        Html
	HtmlFiles   []fidi.File
	hasIndex    bool
	FullNav     bool
	gen         bool
	index       fidi.File
	Assets      []Asset
	Children    []*Page
	Nav         []map[string]any
	Breadcrumbs []map[string]any
	tmpl        *template.Template
	profile     string
}

type BuildOpt func(p *Page)

func NewPage(dir fidi.Tree) *Page {
	page := Page{
		Tree:    dir,
		Css:     GetCss("global"),
		Scripts: GetScripts("global"),
		Html:    GetHtml("global"),
		profile: "global",
		Color:   viper.GetStringMapString("color"),
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

func (pg *Page) Build(opts ...BuildOpt) {
	fmt.Printf("building %s\n", pg.Info().Name)
	for _, opt := range opts {
		opt(pg)
	}
	pg.Render()
}

func (pg *Page) Index() *Page {
	for _, file := range pg.HtmlFiles {
		if file.Base == "index.html" {
			pg.index = file
			pg.hasIndex = true
		}
	}
	return pg
}

func (pg Page) HasIndex() bool {
	for _, file := range pg.HtmlFiles {
		if file.Base == "index.html" {
			return true
		}
	}
	return false
}

func (pg *Page) SetTmpl(tmpl *template.Template) *Page {
	pg.tmpl = tmpl
	return pg
}

func (pg Page) Render() string {
	if pg.gen {
		tmpl := Templates.Lookup("base")
		name := filepath.Join(pg.Info().Path(), "index.html")

		file, err := os.Create(name)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		err = tmpl.Execute(file, pg)
		if err != nil {
			log.Fatal(err)
		}
		return name
	}
	return ""
}

func (pg Page) Content() string {
	var buf bytes.Buffer
	err := pg.tmpl.Execute(&buf, pg)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func (pg *Page) GetChildren() []*Page {
	for _, dir := range pg.Tree.Children() {
		p := NewPage(dir)
		pg.Children = append(pg.Children, p)
	}
	return pg.Children
}

func (pg *Page) SetProfile(pro string) *Page {
	pg.tmpl = GetTemplate(pro)
	pg.profile = pro

	css := GetCss(pro)
	pg.Css = append(pg.Css, css...)

	scripts := GetScripts(pro)
	pg.Scripts = append(pg.Scripts, scripts...)

	html := GetHtml(pro)
	maps.Copy(pg.Html, html)

	mt := pro + ".mime"
	ext := pro + ".ext"
	var items []fidi.File
	switch {
	case viper.IsSet(mt):
		mimes := viper.GetStringSlice(mt)
		items = pg.FilterByMime(mimes...)
	case viper.IsSet(ext):
		exts := viper.GetStringSlice(ext)
		items = pg.FilterByExt(exts...)
	}

	for _, i := range items {
		asset := NewAsset(i, pg.Html)
		pg.Assets = append(pg.Assets, asset)
	}
	return pg
}

func (pg *Page) getBreadcrumbs() *Page {
	var crumbs []map[string]any

	totalP := len(pg.Parents())
	for _, parent := range pg.Parents() {
		totalP--

		path := ".." + strings.Repeat("/..", totalP)
		path = filepath.Join(path, "index.html")

		name := parent.Info().Base
		if parent.Info().Rel() == "." {
			name = "Home"
		}

		link := map[string]any{
			"href":   path,
			"text":   name,
			"indent": parent.Info().Depth,
		}
		crumbs = append(crumbs, link)
	}
	pg.Breadcrumbs = crumbs

	return pg
}

func (pg *Page) getFiles(rel string) []map[string]any {
	var files []map[string]any
	for _, file := range pg.Leaves() {
		if base := file.Base; base != "index.html" {
			url := map[string]any{
				"href":   filepath.Join(rel, base),
				"text":   base,
				"indent": file.Depth,
			}
			files = append(files, url)
		}
	}
	return files
}

func (pg *Page) getNav() *Page {
	var depth []int
	var nav []map[string]any
	for _, p := range pg.Children {
		self := pg.Info().Rel()
		child := p.Info().Rel()

		rel, err := filepath.Rel(self, child)
		if err != nil {
			log.Fatal(err)
		}

		depth = append(depth, p.Info().Depth)

		url := map[string]any{
			"href":   filepath.Join(rel, "index.html"),
			"text":   p.Title,
			"indent": p.Info().Depth,
		}

		if pg.FullNav {
			url["children"] = p.getFiles(rel)
		}

		nav = append(nav, url)
	}

	for idx, d := range lo.Uniq(depth) {
		for _, n := range nav {
			if n["indent"].(int) == d {
				n["indent"] = idx
			}
			if pg.FullNav {
				files := n["children"].([]map[string]any)
				if len(files) > 0 {
					for _, f := range files {
						if f["indent"].(int) == d {
							f["indent"] = idx + 1
						}
					}
				}
			}
		}
	}

	pg.Nav = nav

	return pg
}

func (pg *Page) FilterByExt(ext ...string) []fidi.File {
	return pg.Filter(fidi.ExtFilter(ext...))
}

func (pg *Page) FilterByMime(mime ...string) []fidi.File {
	return pg.Filter(fidi.MimeFilter(mime...))
}
