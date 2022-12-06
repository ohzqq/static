package page

import (
	"bytes"
	"log"
	"path/filepath"
	"static/config"
	"static/files"
)

type Collection struct {
	*Page
	Root     string
	Filetree string
}

func NewCollection(root string, collection config.Collection) Collection {
	col := Collection{
		Root: root,
		Page: NewPageWithChildren(root, collection),
	}
	col.Url = "./index.html"
	for _, p := range col.Children {
		RelativeUrls(root, p)
	}
	col.Filetree = col.Tree()
	return col
}

func (c Collection) Tree() string {
	return c.Content()
}

func (c Collection) Content() string {
	var buf bytes.Buffer
	err := Templates.ExecuteTemplate(&buf, "filetree", c)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (c Collection) Render() []byte {
	var buf bytes.Buffer
	err := Templates.ExecuteTemplate(&buf, "base", c)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func RelativeUrls(root string, pages ...*Page) []*Page {
	for _, page := range pages {
		rel, err := filepath.Rel(root, page.Path)
		if err != nil {
			log.Fatal(err)
		}
		u := filepath.Join(rel, "index.html")
		page.SetUrl("./" + u)

		if page.HasChildren() {
			page.Children = RelativeUrls(root, page.Children...)
		}
	}
	return pages
}

func NewCollectionWithExt(root string, ext ...string) *Collection {
	page := Collection{Root: root}
	page.Page = MakeIndexWithExt(root, ext...)
	page.Files = append(page.Files, files.GlobExt(page.Path, ext...)...)

	return &page
}
