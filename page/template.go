package page

import (
	"idx"
	"idx/config"
	"log"
	"path/filepath"
	"text/template"
)

var (
	TmplFuncs = template.FuncMap{
		"colors": config.RenderColor,
		"color":  config.Colors,
		"Batch":  Batch,
		"Rel":    RelPath,
	}
)

var Templates = template.Must(template.New("").Funcs(TmplFuncs).ParseFS(idx.Static, "static/templates/*"))

func Batch(og []string) [][]string {
	var (
		batch = 4
		files [][]string
	)
	for i := 0; i < len(og); i += batch {
		j := i + batch
		if j > len(og) {
			j = len(og)
		}

		files = append(files, og[i:j])
	}
	return files
}

func RelPath(b, p string) string {
	path, err := filepath.Rel(b, p)
	if err != nil {
		log.Fatal(err)
	}
	return path
}
