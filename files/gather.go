package files

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"
)

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

func GetDirEntries(name string) []os.DirEntry {
	entries, err := os.ReadDir(name)
	if err != nil {
		log.Fatal(err)
	}
	return entries
}
