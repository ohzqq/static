package page

import (
	"embed"
	"idxgen/config"
	"text/template"
)

var (
	//go:embed templates/*
	templateDir embed.FS

	TmplFuncs = template.FuncMap{
		"colors": config.RenderColor,
		"color":  config.Colors,
	}
)

var Templates = template.Must(template.New("").Funcs(TmplFuncs).ParseFS(templateDir, "templates/*"))
