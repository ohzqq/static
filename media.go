package static

import (
	"bytes"
	"encoding/base64"
	"log"
	"mime"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Gallery struct {
	Meta  Meta     `json:"meta,omitempty"`
	Media []*Media `json:"media"`
}

type Meta struct {
	Title   string   `json:"title"`
	Creator string   `json:"creator"`
	Subject []string `json:"subject"`
}

type Media struct {
	Img       string   `json:"img,omitempty"`
	Video     string   `json:"video,omitempty"`
	Caption   string   `json:"caption,omitempty"`
	Thumbnail string   `json:"thumbnail,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

func Thumb(input string, output ...string) Media {
	var thumb []byte
	var media Media

	ext := filepath.Ext(input)
	mt := mime.TypeByExtension(ext)

	switch {
	case strings.Contains(mt, "video"):
		thumb = VideoThumb(input)
		media.Video = input
	case strings.Contains(mt, "image"):
		thumb = ImageThumb(input)
		media.Img = input
	}
	media.Thumbnail = ThumbToBase64(thumb)

	return media
}

func VideoThumb(input string) []byte {
	inArgs := ffmpeg.KwArgs{
		"y":           "",
		"loglevel":    "error",
		"hide_banner": "",
	}
	outArgs := ffmpeg.KwArgs{
		"c:v":      "mjpeg",
		"frames:v": 1,
		"f":        "image2",
	}

	ff := ffmpeg.Input(input, inArgs).
		Filter("thumbnail", ffmpeg.Args{}).
		Filter("scale", ffmpeg.Args{"w=268:h=-2:flags=lanczos"}).
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

func ImageThumb(path string) []byte {
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
