package page

import (
	"idx/config"
	"idx/files"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

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

func MakeIndexWithMime(root string, col config.Collection) Page {
	idx := Page{Collection: col}
	idx.Path = filepath.Join(idx.Path, root)
	entries := files.GetDirEntries(idx.Path)

	for _, e := range entries {
		var child Page
		fp := filepath.Join(idx.Path, e.Name())
		if e.IsDir() {
			child = MakeIndexWithMime(fp, col)
			child.Files = append(child.Files, files.GlobMime(fp, col.Mime)...)
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
