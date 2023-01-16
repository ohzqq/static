package static

import "github.com/ohzqq/fidi"

type Page struct {
	fidi.Tree
	Index    fidi.File
	HasIndex bool
	Assets   []fidi.File
	//Children []fidi.File
	Nav []fidi.Dir
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
		}
	}

	return page
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
