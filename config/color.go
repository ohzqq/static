package config

import (
	"bytes"
	"log"
	"text/template"
)

type Color struct {
	Bg     string
	Fg     string
	White  string
	Black  string
	Grey   string
	Yellow string
	Red    string
	Pink   string
	Cyan   string
	Blue   string
	Green  string
	Purple string
}

var ColorTmpl = template.Must(template.New("cssColor").Parse(CssColorTmpl))

func RenderColor() string {
	var buf bytes.Buffer
	err := ColorTmpl.ExecuteTemplate(&buf, "cssColor", Opts.Color)
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

const CssColorTmpl = `:root {
	{{with .Bg}}--bg: {{.}};{{end}}
	{{with .Fg}}--fg: {{.}};{{end}}
	{{with .White}}--white: {{.}};{{end}}
	{{with .Black}}--black: {{.}};{{end}}
	{{with .Grey}}--grey: {{.}};{{end}}
	{{with .Yellow}}--yellow: {{.}};{{end}}
	{{with .Red}}--red: {{.}};{{end}}
	{{with .Pink}}--pink: {{.}};{{end}}
	{{with .Cyan}}--cyan: {{.}};{{end}}
	{{with .Blue}}--blue: {{.}};{{end}}
	{{with .Green}}--green: {{.}};{{end}}
	{{with .Purple}}--purple: {{.}};{{end}}
	--darken: brightness(90%);
}
`
