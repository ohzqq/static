package page

import (
	"bytes"
	"html/template"
	"idxgen/config"
	"idxgen/files"
	"log"
	"os"
	"path/filepath"
)

type Meta struct {
	Title string   `toml:"title"`
	Tags  []string `toml:"tags"`
}

type Collection struct {
	Page
	Root string
}

type Page struct {
	Meta     Meta
	Body     string
	Type     string
	Path     string `toml:"path"`
	Files    []string
	Children []Page
	config.Collection
}

func NewCollection(root string, ext ...string) *Collection {
	page := Collection{Root: root}
	page.Page = MakeIndexWithExt(root, ext...)
	page.Files = append(page.Files, files.GlobExt(page.Path, ext...)...)

	return &page
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

func (p Page) Parse() []byte {
	t := template.Must(template.New("imagePage").ParseFiles(p.Template))

	var buf bytes.Buffer
	err := t.Execute(&buf, p)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func DumpFile(path string) []byte {
	t, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
