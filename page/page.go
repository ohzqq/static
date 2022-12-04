package page

import (
	"bytes"
	"html/template"
	"idxgen/config"
	"idxgen/files"
	"log"
	"path/filepath"
)

type Meta struct {
	Title string   `toml:"title"`
	Tags  []string `toml:"tags"`
}

type Page struct {
	Meta     Meta
	Type     string
	Path     string `toml:"path"`
	Files    []string
	Children []Page
	config.Collection
}

func NewPage(root, collection string) Page {
	page := Page{
		Path:       root,
		Collection: config.GetCollection(collection),
	}
	page.Files = append(page.Files, files.GlobMime(page.Path, page.Mime)...)

	return page
}

func NewPageWithChildren(root, col string) Page {
	collection := config.GetCollection(col)
	page := MakeIndexWithMime(root, collection.Mime)
	page.Collection = collection
	page.Files = append(page.Files, files.GlobMime(page.Path, page.Mime)...)

	return page
}

func (p Page) Title() string {
	return filepath.Base(p.Path)
}

func (p Page) HasChildren() bool {
	return len(p.Children) > 0
}

func (p Page) Render() string {
	var buf bytes.Buffer
	err := Templates.ExecuteTemplate(&buf, "base", p)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (p Page) Content() string {
	var buf bytes.Buffer

	if p.Template != "" {
		t := template.Must(template.New("content").ParseFiles(p.Template))
		err := t.ExecuteTemplate(&buf, "content", p)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := Templates.ExecuteTemplate(&buf, "body"+p.Mime, p)
		if err != nil {
			log.Fatal(err)
		}
	}

	return buf.String()
}
