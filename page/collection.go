package page

import (
	"idx/config"
	"idx/files"
)

type Collection struct {
	Page
	Root     string
	Filetree string
}

func NewCollection(root string, collection config.Collection) Collection {
	page := Collection{
		Root: root,
		Page: NewPageWithChildren(root, collection),
	}
	page.Filetree = page.Tree()

	return page
}

//func (c Collection) MakePages() {
//  if c.HasChildren() {
//    for _, page := range c.Children {
//    }
//  }
//}

func NewCollectionWithExt(root string, ext ...string) *Collection {
	page := Collection{Root: root}
	page.Page = MakeIndexWithExt(root, ext...)
	page.Files = append(page.Files, files.GlobExt(page.Path, ext...)...)

	return &page
}
