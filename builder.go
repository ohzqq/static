package static

import (
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ohzqq/fidi"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

type Builder struct {
	Nav          bool
	FullNav      bool
	Gen          bool
	Regen        bool
	isCollection bool
	Tmpl         *template.Template
	Profile      string
	Input        string
}

func New(path string) *Builder {
	return &Builder{
		Input: path,
	}
}

func (b *Builder) Collection() *Builder {
	b.isCollection = true
	b.Nav = true
	return b
}

func (b *Builder) Page() *Builder {
	b.isCollection = false
	return b
}

func (b Builder) Opts() []BuildOpt {
	var opts []BuildOpt

	if b.isCollection {
		opts = append(opts, Collection())
	}

	if b.Nav {
		opts = append(opts, Nav(b.FullNav))
	}

	switch {
	case b.Gen:
		opts = append(opts, Gen())
	case b.Regen:
		opts = append(opts, Regen())
	}

	if b.Profile != "" {
		opts = append(opts, Profile(b.Profile))
	}

	return opts
}

func (b *Builder) Build() {
	if b.Input == "" {
		log.Fatal("no input")
	}
	tree := fidi.NewTree(b.Input)

	page := NewPage(tree)
	page.Build(b.Opts()...)

	if b.isCollection {
		for _, child := range page.Children {
			child.Build(b.Opts()...)
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
		p.FullNav = full
		p.getBreadcrumbs()
		p.getNav()
	}
}

func Collection() BuildOpt {
	return func(page *Page) {
		for _, dir := range page.Tree.Children() {
			p := NewPage(dir)
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
func getBreadcrumbs(tree fidi.Tree) []map[string]any {
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

	return crumbs
}

func getFiles(page *Page, rel string) []map[string]any {
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

func getNav(page *Page) []map[string]any {
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
			url["children"] = getFiles(p, rel)
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

	return nav
}
