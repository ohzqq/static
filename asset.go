package static

import (
	"bytes"
	"encoding/base64"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/disintegration/imaging"
	"github.com/ohzqq/fidi"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Asset struct {
	fidi.File
	Attributes map[string]any
	Html       Html
	NoThumbs   bool
	Tag        string
}

func NewAsset(file fidi.File, noThumbs bool, tags ...Html) Asset {
	var html Html
	if len(tags) > 0 {
		html = tags[0]
	}
	asset := Asset{
		File:       file,
		Html:       html,
		Attributes: make(map[string]any),
	}
	return asset
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
		if at, ok := a.Html[a.Tag]; ok {
			a.Attributes = at
			a.Attributes["src"] = a.Base
		}
	case a.IsImage():
		a.Tag = "img"
		if at, ok := a.Html[a.Tag]; ok {
			a.Attributes = at
		}
		a.Attributes["src"] = Thumbnail(a.Path())
		a.Attributes["alt"] = a.Base
	case a.IsVideo():
		a.Tag = "video"
		if at, ok := a.Html[a.Tag]; ok {
			a.Attributes = at
		}
		a.Attributes["src"] = a.Base
		a.Attributes["poster"] = ExtractThumbFromVideo(a.File)
	}

	var buf bytes.Buffer
	err := assetTmpl.Execute(&buf, a)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func ExtractThumbFromVideo(file fidi.File) string {
	out := file.Copy().Ext(".jpg").Prefix("thumb-")
	tmp := filepath.Base(out.String())
	tmp = filepath.Join(os.TempDir(), tmp)
	defer os.Remove(tmp)

	ff := ffmpeg.Input(file.Path(), ffmpeg.KwArgs{"y": ""}).
		Filter("thumbnail", ffmpeg.Args{}).
		Output(tmp, ffmpeg.KwArgs{"frames:v": 1})

	args := ff.GetArgs()
	cmd := exec.Command("ffmpeg", args...)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	base := Thumbnail(tmp)

	return base
}

func Thumbnail(path string) string {
	src, err := imaging.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	thumb := imaging.Fit(src, 268, 150, imaging.Lanczos)

	var buf bytes.Buffer
	err = imaging.Encode(&buf, thumb, imaging.JPEG)
	if err != nil {
		log.Fatal(err)
	}

	base := "data:image/jpeg;base64,"
	base += base64.StdEncoding.EncodeToString(buf.Bytes())

	return base
}

var assetTmpl = template.Must(template.New("asset").Parse(
	`<{{.Tag}}
	{{- range $key, $val := .Attributes}} 
	{{- if ne $key "autoplay"}}
	{{- if ne $key "controls"}}
	{{$key}}="{{$val}}"
	{{- end -}}
	{{- end -}}
	{{- end -}}
	></{{.Tag}}>
`))
