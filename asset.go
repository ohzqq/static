package static

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/ohzqq/fidi"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/yaml.v2"
)

type Assets []Asset

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
			a.Attributes["poster"] = Thumbnail(a.File, "base64")
		}
	case "img":
		a.Attributes["data-original"] = a.Base
		if !noThumbs() {
			a.Attributes["src"] = Thumbnail(a.File, "base64")
		}
	}

	return a
}

func (as Assets) Export(format, path string) {
	var data []byte
	switch format {
	case "yaml", "yml":
		data = as.toYaml()
	case "json":
		data = as.toJson()
	}

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func (as Assets) toJson() []byte {
	data, err := json.Marshal(as)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func (as Assets) toYaml() []byte {
	data, err := yaml.Marshal(as)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func Thumbnail(input fidi.File, output ...string) string {
	var thumb []byte
	switch {
	case strings.Contains(input.Mime, "video"):
		thumb = videoThumb(input.Path())
	case strings.Contains(input.Mime, "image"):
		thumb = imageThumb(input.Path())
	}

	name := input.Copy().Ext(".jpg").Prefix("thumb-").String()

	if len(output) > 0 {
		name = output[0]
	}

	switch name {
	case "base64":
		return ThumbToBase64(thumb)
	default:
		out := filepath.Join(input.Dir, name)
		file, err := os.Create(out)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = file.Write(thumb)
		if err != nil {
			log.Fatal(err)
		}
		return name
	}
}

func videoThumb(input string) []byte {
	inArgs := ffmpeg.KwArgs{
		"y":           "",
		"loglevel":    "quiet",
		"hide_banner": "",
	}
	outArgs := ffmpeg.KwArgs{
		"c:v":      "mjpeg",
		"frames:v": 1,
		"f":        "image2",
	}

	ff := ffmpeg.Input(input, inArgs).
		Filter("thumbnail", ffmpeg.Args{}).
		Output("pipe:1", outArgs)

	args := ff.GetArgs()
	cmd := exec.Command("ffmpeg", args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	return out.Bytes()
}

func imageThumb(path string) []byte {
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

	return buf.Bytes()
}

func ThumbToBase64(data []byte) string {
	base := "data:image/jpeg;base64,"
	base += base64.StdEncoding.EncodeToString(data)
	return base
}
