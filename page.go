package static

import "github.com/ohzqq/fidi"

type Page struct {
	fidi.Dir
}

func NewPage(dir fidi.Dir) Page {
	return Page{
		Dir: dir,
	}
}

func (p Page) Children() []fidi.File {
	return p.Files
}

func (p Page) FilterByExt(ext ...string) []fidi.File {
	return p.Filter(fidi.ExtFilter(ext...))
}

func (p Page) FilterByMime(mime ...string) []fidi.File {
	return p.Filter(fidi.MimeFilter(mime...))
}
