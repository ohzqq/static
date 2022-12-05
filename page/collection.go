package page

import "idxgen/files"

type Collection struct {
	Page
	Root string
}

func NewCollection(root, col string) Collection {
	page := Collection{
		Root: root,
		Page: NewPageWithChildren(root, col),
	}

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
