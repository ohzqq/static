package static

import (
	"embed"
	"html/template"
	"static/config"
)

var (
	//go:embed static/*
	Static embed.FS
)

var (
	TmplFuncs = template.FuncMap{
		"color": config.Colors,
		"Batch": Batch,
	}
)

var Templates = template.Must(template.New("").Funcs(TmplFuncs).ParseFS(static.Static, "static/templates/*"))

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
