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

func NewPage(dir fidi.Tree) *Page {
	page := Page{
		Tree:    dir,
		Css:     GetCss("global"),
		Scripts: GetScripts("global"),
		Html:    GetHtml("global"),
		profile: "global",
		Color:   viper.GetStringMapString("color"),
		Opts:    &Builder{},
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

func (p *Page) Build(opts ...BuildOpt) {
	fmt.Printf("building %s\n", p.Info().Name)
	for _, opt := range opts {
		opt(p)
	}
	p.Render()
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
		p := NewPage(dir)
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
	url["indent"] = p.Info().Depth
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

func (tree *Page) getBreadcrumbs() *Page {
	var crumbs []map[string]any

	totalP := len(tree.Parents())
	for _, parent := range tree.Parents() {
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
	tree.Breadcrumbs = crumbs

	return tree
}

func (page *Page) getFiles(rel string) []map[string]any {
	var files []map[string]any
	for _, file := range page.Leaves() {
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

func (page *Page) getNav() *Page {
	var depth []int
	var nav []map[string]any
	for _, p := range page.Children {
		self := page.Info().Rel()
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

		if page.FullNav {
			url["children"] = p.getFiles(rel)
		}

		nav = append(nav, url)
	}

	for idx, d := range lo.Uniq(depth) {
		for _, n := range nav {
			if n["indent"].(int) == d {
				n["indent"] = idx
			}
			if page.FullNav {
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

	page.Nav = nav

	return page
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
