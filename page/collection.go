package page

import "idxgen/files"

type Collection struct {
	Page
	Root string
}

func NewCollection(root string, ext ...string) *Collection {
	page := Collection{Root: root}
	page.Page = MakeIndexWithExt(root, ext...)
	page.Files = append(page.Files, files.GlobExt(page.Path, ext...)...)

	return &page
}
