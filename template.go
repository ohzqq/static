package static

import (
	"path/filepath"
	"text/template"
)

var (
	TmplFuncs = template.FuncMap{
		"color": Colors,
		"Batch": Batch,
		"Dir":   filepath.Dir,
	}
)

var Templates = template.Must(template.New("").Funcs(TmplFuncs).ParseFS(Static, "static/templates/*"))

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
