package static

import "github.com/ohzqq/fidi"

type Page struct {
	fidi.Dir
	Children []fidi.File
	Nav      []fidi.Dir
}

func NewPage(dir fidi.Dir) Page {
	return Page{
		Dir: dir,
	}
}

func (p *Page) FilterByExt(ext ...string) *Page {
	p.Children = p.Filter(fidi.ExtFilter(ext...))
	return p
}

func (p *Page) FilterByMime(mime ...string) *Page {
	p.Children = p.Filter(fidi.MimeFilter(mime...))
	return p
}
