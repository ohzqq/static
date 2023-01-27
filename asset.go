package static

import (
	"bytes"
	"encoding/base64"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/ohzqq/fidi"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Assets struct {
	Assets []Assets `json:"Assets" yaml:"Assets"`
}

type Asset struct {
	fidi.File  `json:"-" yaml:"-"`
	Tag        string         `json:"Tag" yaml:"Tag"`
	Attributes map[string]any `json:"Attributes" yaml:"Attributes"`
}

func NewAsset(file fidi.File, tags ...Html) Asset {
	html := getHtml()
	a := Asset{
		File:       file,
		Attributes: make(map[string]any),
	}

	switch {
	case strings.Contains(a.Mime, "audio"):
		a.Tag = "audio"
	case strings.Contains(a.Mime, "video"):
		a.Tag = "video"
	case strings.Contains(a.Mime, "image"):
		a.Tag = "img"
	}

	if at, ok := html[a.Tag]; ok {
		a.Attributes = at
	}
	a.Attributes["src"] = a.Base
	a.Attributes["mime"] = a.Mime
	a.Attributes["caption"] = ""

	switch a.Tag {
	case "video":
		a.Attributes["poster"] = a.Base
		if !noThumbs() {
			a.Attributes["poster"] = ExtractThumbFromVideo(a.File)
		}
	case "img":
		a.Attributes["data-original"] = a.Base
		if !noThumbs() {
			a.Attributes["src"] = Thumbnail(a.Path())
		}
	}

	return a
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
