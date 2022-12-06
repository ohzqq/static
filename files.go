package static

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

func Write(path string, page []byte) error {
	out := filepath.Join(path, "index.html")

	err := os.WriteFile(out, page, 0666)
	if err != nil {
		return fmt.Errorf("Rendering %s failed with error %s\n", out, err)
	}
	fmt.Printf("Rendered %s\n", out)
	return nil
}

func GlobMime(path, mtype string) []string {
	var files []string
	for _, entry := range GetDirEntries(path) {
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
	entries := GetDirEntries(path)
	for _, entry := range entries {
		ePath := filepath.Join(entry.Name())
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
