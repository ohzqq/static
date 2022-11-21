package files

import (
	"log"
	"os"
	"path/filepath"
)

type Site struct {
	Index
	Root string
}

func (s *Site) GetPages(root string) *Site {
	s.Root = root
	s.Index = MakeIndex(root)
	return s
}

type Index struct {
	Body     string
	Meta     string
	Children []Index
	Path     string `toml:"path"`
}

func MakeIndex(root string) Index {
	var idx Index
	idx.Path = filepath.Join(idx.Path, root)
	entries := GetDirEntries(idx.Path)
	for _, e := range entries {
		if e.IsDir() {
			child := MakeIndex(filepath.Join(idx.Path, e.Name()))
			idx.Children = append(idx.Children, child)
		}
		switch name := e.Name(); name {
		case "body.html":
			idx.Body = name
		case "meta.toml":
			idx.Meta = name
		}
	}
	return idx
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
