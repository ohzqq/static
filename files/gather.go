package files

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"
)

type Page struct {
	Index
	Root string
}

type Index struct {
	Type     string
	Body     string
	Meta     string
	Files    []string
	Children []Index
	Path     string `toml:"path"`
}

func MakePage(root string, ext ...string) *Page {
	page := Page{Root: root}
	page.Index = MakeIndex(root, ext...)
	page.Files = append(page.Files, GlobExt(page.Path, ext...)...)

	return &page
}

func MakeIndex(root string, ext ...string) Index {
	var idx Index
	idx.Path = filepath.Join(idx.Path, root)
	entries := GetDirEntries(idx.Path)

	for _, e := range entries {
		var child Index
		fp := filepath.Join(idx.Path, e.Name())
		if e.IsDir() {
			child = MakeIndex(fp, ext...)
			child.Files = append(child.Files, GlobExt(fp, ext...)...)
			idx.Children = append(idx.Children, child)
		}
	}
	return idx
}

func GlobExt(path string, ext ...string) []string {
	var files []string
	for _, entry := range GetDirEntries(path) {
		ePath := filepath.Join(path, entry.Name())
		if eExt := filepath.Ext(ePath); slices.Contains(ext, eExt) {
			files = append(files, ePath)
		}
	}
	return files
}

type Meta struct {
	Title    string   `toml:"title"`
	Template string   `toml:"template"`
	Tags     []string `toml:"tags"`
	Cmd      `toml:"cmd"`
}

type Cmd struct {
	Bin  string   `toml:"bin"`
	Args []string `toml:"args"`
}

func GetDirEntries(name string) []os.DirEntry {
	//abs, err := filepath.Abs(name)
	//if err != nil {
	//  log.Fatal(err)
	//}
	entries, err := os.ReadDir(name)
	if err != nil {
		log.Fatal(err)
	}
	return entries
}
