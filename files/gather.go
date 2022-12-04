package files

import (
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

func GlobMime(path, mtype string) []string {
	var files []string
	for _, entry := range GetDirEntries(path) {
		//ePath := filepath.Join(path, entry.Name())
		ePath := filepath.Join(entry.Name())
		eExt := filepath.Ext(ePath)
		mt := mime.TypeByExtension(eExt)
		if strings.Contains(mt, mtype) {
			files = append(files, ePath)
		}
	}
	return files
}

func GlobExt(path string, ext ...string) []string {
	var files []string
	for _, entry := range GetDirEntries(path) {
		ePath := filepath.Join(path, entry.Name())
		eExt := filepath.Ext(ePath)
		if slices.Contains(ext, eExt) {
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
