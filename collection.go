package static

import (
	"bytes"
	"log"
	"strings"
	"text/template"

	"github.com/ohzqq/fidi"
)

type Collection struct {
	*Page
}

func NewCollection(path string, profile ...string) Collection {
	tree := fidi.NewTree(path)

	col := Collection{
		Page: NewPage(tree),
	}
	col.root = path

	if len(profile) > 0 {
		col.Profile(profile[0])
	}

	col.GetChildren()

	for _, page := range col.Children {
		page.GetChildren()
	}

	return col
}

type Asset struct {
	fidi.File
	Attributes map[string]any
	Html       Html
	Tag        string
}

func (a Asset) IsAudio() bool {
	return strings.Contains(a.Mime, "audio")
}

func (a Asset) IsVideo() bool {
	return strings.Contains(a.Mime, "video")
}

func (a Asset) IsImage() bool {
	return strings.Contains(a.Mime, "image")
}

func (a *Asset) Render() string {
	switch {
	case a.IsAudio():
		a.Tag = "audio"
	case a.IsImage():
		a.Tag = "img"
	case a.IsVideo():
		a.Tag = "video"
	}

	if at, ok := a.Html[a.Tag]; ok {
		a.Attributes = at
	}

	var buf bytes.Buffer
	err := assetTmpl.Execute(&buf, a)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

var assetTmpl = template.Must(template.New("asset").Parse(`<{{.Tag}}{{range $key, $val := .Attributes}} {{$key}}="{{$val}}"{{end}} src="{{.Base}}"></{{.Tag}}>`))
