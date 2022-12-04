package page

import (
	"idxgen/files"
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

func MakeIndexWithMime(root string, mime string) Page {
	var idx Page
	idx.Path = filepath.Join(idx.Path, root)
	entries := files.GetDirEntries(idx.Path)

	for _, e := range entries {
		var child Page
		fp := filepath.Join(idx.Path, e.Name())
		if e.IsDir() {
			child = MakeIndexWithMime(fp, mime)
			child.Files = append(child.Files, files.GlobMime(fp, mime)...)
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
