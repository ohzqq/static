package page

import (
	"idxgen/config"
	"idxgen/files"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
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

func MakeIndexWithExt(root string, ext ...string) Page {
	var idx Page
	idx.Path = filepath.Join(idx.Path, root)
	entries := files.GetDirEntries(idx.Path)

	for _, e := range entries {
		var child Page
		fp := filepath.Join(idx.Path, e.Name())
		if e.IsDir() {
			child = MakeIndexWithExt(fp, ext...)
			child.Files = append(child.Files, files.GlobExt(fp, ext...)...)
			idx.Children = append(idx.Children, child)
		}
		switch name := e.Name(); name {
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

func MakeIndexWithExt(root string, ext ...string) Page {
	var idx Page
	idx.Path = filepath.Join(idx.Path, root)
	entries := files.GetDirEntries(idx.Path)

	for _, e := range entries {
		var child Page
		fp := filepath.Join(idx.Path, e.Name())
		if e.IsDir() {
			child = MakeIndexWithExt(fp, ext...)
			child.Files = append(child.Files, files.GlobExt(fp, ext...)...)
			idx.Children = append(idx.Children, child)
		}
		switch name := e.Name(); name {
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

func NewPage(root, collection string) Page {
	page := Page{
		Path:       root,
		Collection: config.GetCollection(collection),
	}
	page.Files = append(page.Files, files.GlobMime(page.Path, page.Mime)...)

	return page
}

func NewPageWithChildren(root, collection string) Page {
	page := MakeIndexWithExt(root)
	page := Page{
		Path:       root,
		Collection: config.GetCollection(collection),
	}
	page.Files = append(page.Files, files.GlobMime(page.Path, page.Mime)...)

	return page
}

func (p Page) Title() string {
	return filepath.Base(p.Path)
}
