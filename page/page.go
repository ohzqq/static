package page

import (
	"idxgen/config"
	"idxgen/files"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Collection struct {
	Page
	Root string
}

type Page struct {
	Meta
	Body     string
	Type     string
	Path     string `toml:"path"`
	Files    []string
	Children []Page
	config.Collection
}

type Meta struct {
	Title string
	Tags  []string
}

func New(root string, ext ...string) *Collection {
	page := Collection{Root: root}
	page.Page = MakeIndex(root, ext...)
	page.Files = append(page.Files, files.GlobExt(page.Path, ext...)...)

	return &page
}

func MakeIndex(root string, ext ...string) Page {
	var idx Page
	idx.Path = filepath.Join(idx.Path, root)
	entries := files.GetDirEntries(idx.Path)

	for _, e := range entries {
		var child Page
		fp := filepath.Join(idx.Path, e.Name())
		if e.IsDir() {
			child = MakeIndex(fp, ext...)
			child.Files = append(child.Files, files.GlobExt(fp, ext...)...)
			idx.Children = append(idx.Children, child)
		}
		switch name := e.Name(); name {
		case "body.html":
			idx.Body = name
		case "meta.toml":
			t, err := os.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}
			toml.Unmarshal(t, &idx.Meta)
		}
	}
	return idx
}

func (p Page) Title() string {
	return filepath.Base(p.Path)
}
