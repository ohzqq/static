package static

import (
	"fmt"
	"path/filepath"
	"text/template"
)

var (
	TmplFuncs = template.FuncMap{
		"Batch": Batch,
		"Dir":   filepath.Dir,
	}
)

var publicTemplates = template.Must(template.New("").Funcs(TmplFuncs).ParseFS(Public, "static/templates/*"))

func GetTemplates() []string {
	var tmpl []string
	for _, pro := range Profiles() {
	}
	fmt.Printf("profiles %s\n", tmpl)
	return tmpl
}

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
