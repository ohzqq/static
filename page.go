package static

import (
	"io/fs"
	"log"
	"strings"

	"github.com/ohzqq/fidi"
)

type Page struct {
	fidi.Tree
	Profile
	HasIndex bool
	Index    fidi.File
	Url      string
	Assets   []fidi.File
	Nav      []fidi.Dir
}

func NewPage(dir fidi.Tree) Page {
	page := Page{
		Tree: dir,
	}

	files := page.Filter(fidi.ExtFilter(".html"))
	for _, file := range files {
		if file.Base == "index.html" {
			page.HasIndex = true
			page.Index = file
			page.Url = "./" + file.Rel()
		}
	}

	return page
}

func (p *Page) SetProfile(pro string) *Page {
	p.Profile = GetProfile(pro)
	return p
}

func (p *Page) CreateIndex() *Page {
	if !p.HasIndex {
		//Create filetree page
	}
	return p
}

func (p *Page) FilterByExt(ext ...string) *Page {
	p.Assets = p.Filter(fidi.ExtFilter(ext...))
	return p
}

func (p *Page) FilterByMime(mime ...string) *Page {
	p.Assets = p.Filter(fidi.MimeFilter(mime...))
	return p
}

func (p Page) ReadCss() []string {
	var assets []string
	for _, css := range p.Css {
		var f fs.FS
		if strings.HasPrefix(css, "static") {
			f = Public
		} else {
			f = UserCfg
		}
		d, err := fs.ReadFile(f, css)
		if err != nil {
			log.Fatal(err)
		}
		assets = append(assets, string(d))
	}
	return assets
}

func (p Page) ReadScripts() []string {
	var assets []string
	for _, script := range p.Scripts {
		var f fs.FS
		if strings.HasPrefix(script, "static") {
			f = Public
		} else {
			f = UserCfg
		}
		d, err := fs.ReadFile(f, script)
		if err != nil {
			log.Fatal(err)
		}
		assets = append(assets, string(d))
	}
	return assets
}
