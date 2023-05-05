package static

import (
	"encoding/json"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/ohzqq/fidi"
	"gopkg.in/yaml.v2"
)

type Assets []Asset

type mimeType string

func MimeType(path string) mimeType {
	ext := filepath.Ext(path)
	return mimeType(mime.TypeByExtension(ext))
}

func (mt mimeType) IsMedia() bool {
	if mt.IsVideo() || mt.IsImage() {
		return true
	}
	return false
}

func (mt mimeType) IsVideo() bool {
	if strings.Contains(string(mt), "video") {
		return true
	}
	return false
}

func (mt mimeType) IsImage() bool {
	if strings.Contains(string(mt), "image") {
		return true
	}
	return false
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
		thumb = VideoThumb(input.Path())
	case strings.Contains(input.Mime, "image"):
		thumb = ImageThumb(input.Path())
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
