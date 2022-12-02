package page

import (
	"idxgen/config"
	"idxgen/files"
	"path/filepath"
)

type Page struct {
	Index
	Root string
}

type Index struct {
	Type       string
	Collection config.Collection
	Files      []string
	Children   []Index
	Path       string `toml:"path"`
}

func New(root string, ext ...string) *Page {
	page := Page{Root: root}
	page.Index = MakeIndex(root, ext...)
	page.Files = append(page.Files, files.GlobExt(page.Path, ext...)...)

	return &page
}

func MakeIndex(root string, ext ...string) Index {
	var idx Index
	idx.Path = filepath.Join(idx.Path, root)
	entries := files.GetDirEntries(idx.Path)

	for _, e := range entries {
		var child Index
		fp := filepath.Join(idx.Path, e.Name())
		if e.IsDir() {
			child = MakeIndex(fp, ext...)
			child.Files = append(child.Files, files.GlobExt(fp, ext...)...)
			idx.Children = append(idx.Children, child)
		}
	}
	return idx
}
