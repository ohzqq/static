package static

import (
	"bytes"
	"log"
	"path/filepath"
)

type Collection struct {
	*Page
	Root     string
	Filetree string
}

func NewCollection(root string) Collection {
	col := Collection{
		Root: root,
		Page: NewPage(root),
	}
	col.Url = "./index.html"
	col.Filetree = col.Tree()
	return col
}

func (c Collection) Tree() string {
	return c.Content()
}

func (c Collection) Content() string {
	for _, p := range c.Children {
		RelativeUrls(c.Root, p)
	}
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
		println(page.Url)

		if page.HasChildren() {
			page.Children = RelativeUrls(root, page.Children...)
		}
	}
	return pages
}
