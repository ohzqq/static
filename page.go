package static

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ohzqq/fidi"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Page struct {
	fidi.Tree
	Title       string
	Input       string
	css         []string
	scripts     []string
	Color       map[string]string
	Html        Html
	HtmlFiles   []fidi.File
	HasChildren bool
	hasIndex    bool
	FullNav     bool
	index       fidi.File
	Assets      Assets
	filters     []fidi.Filter
	Children    []*Page
	Nav         []map[string]any
	Breadcrumbs []map[string]any
	tmpl        *template.Template
}

type BuildOpt func(p *Page)

func NewPage(dir fidi.Tree) *Page {
	page := Page{
		Tree:  dir,
		Color: viper.GetStringMapString("color"),
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

func (pg *Page) Build() {
	pg.tmpl = getTemplate()
	pg.Html = getHtml()

	if indexOnly() {
		pg.Nav = pg.setFiles(pg.Info().Rel())
		tmpl := Templates.Lookup("filterableList")
		pg.SetTmpl(tmpl)
	}
	//fmt.Printf("index only %v\n", viper.GetString("build.template"))

	pg.GetAssets()
	if viper.GetBool("build.assets") {
		f := viper.GetString("build.format")
		p := pg.Info().Rename("assets").Ext("." + f).String()
		p = filepath.Join(pg.Info().Dir, pg.Info().Base, p)
		pg.Assets.Export(f, p)
	}

	if recurse() {
		pg.setChildren()
		pg.setNav()
		pg.setBreadcrumbs()
	}

	pg.HasChildren = len(pg.Children) > 0

	if !pg.HasIndex() || regen() {
		fmt.Printf("building %s\n", pg.Info().Name)
		pg.Render()
	}
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
	tmpl := Templates.Lookup("base")
	name := filepath.Join(pg.Info().Path(), "index.html")
	//pg.saveContent()

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

func (pg Page) Content() string {
	for _, file := range pg.HtmlFiles {
		if file.Base == "content.html" {
			h, err := os.ReadFile(file.Path())
			if err != nil {
				log.Fatal(err)
			}
			return string(h)
		}
	}

	var buf bytes.Buffer
	err := pg.tmpl.Execute(&buf, pg)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func (pg Page) saveContent() {
	name := pg.Title
	if pg.Title == "Home" {
		name = "index"
	}
	name += ".html"
	name = filepath.Join(pg.Input, name)

	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(pg.Content()))
	if err != nil {
		log.Fatal(err)
	}
}

func (pg Page) Css() []string {
	return parsePageResources(".css")
}

func (pg Page) Scripts() []string {
	return parsePageResources(".scripts")
}

func (pg *Page) setChildren() []*Page {
	for _, dir := range pg.Tree.Children() {
		p := NewPage(dir)
		pg.Children = append(pg.Children, p)
	}
	return pg.Children
}

func (pg *Page) GetAssets() []Asset {
	var filters []fidi.Filter
	m := parseFilterKind(".mime")
	if hasMimes() {
		m = mimes()
	}
	if len(m) > 0 {
		filters = append(filters, fidi.MimeFilter(m...))
	}

	e := parseFilterKind(".ext")
	if hasExts() {
		e = exts()
	}
	if len(e) > 0 {
		filters = append(filters, fidi.ExtFilter(e...))
	}

	items := pg.Filter(filters...)

	var assets []Asset
	for _, i := range items {
		asset := NewAsset(i)
		assets = append(assets, asset)
	}
	pg.Assets = assets

	return assets
}

func (pg *Page) setBreadcrumbs() *Page {
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

func (pg Page) AssetsToJson() string {
	data, err := json.Marshal(pg.Assets)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func (pg Page) AssetsToYaml() string {
	data, err := yaml.Marshal(pg.Assets)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func (pg *Page) setFiles(rel string) []map[string]any {
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

func (pg *Page) setNav() *Page {
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

		if listAll() {
			url["children"] = p.setFiles(rel)
		}

		nav = append(nav, url)
	}

	for idx, d := range lo.Uniq(depth) {
		for _, n := range nav {
			if n["indent"].(int) == d {
				n["indent"] = idx
			}
			if listAll() {
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
