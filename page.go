package static

import (
	"bytes"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
)

type GlobType int

const (
	MimeType GlobType = iota
	Extension
)

type Meta struct {
	Title string   `toml:"title"`
	Tags  []string `toml:"tags"`
}

type Page struct {
	Meta     Meta
	glob     GlobType
	category string
	Category
	Url      string
	Path     string `toml:"path"`
	Files    []string
	Children []*Page
}

func NewPage(root string) *Page {
	page := Page{
		Path:     root,
		Category: DefaultCategory(),
	}
	return &page
}

func (p *Page) SetCategory(c string) *Page {
	if c == "video" || c == "image" {
		c = "swiper"
	}
	p.category = c
	p.Category = GetCategory(c)
	return p
}

func (p *Page) GlobMime(mime ...string) *Page {
	p.glob = MimeType
	if len(mime) > 0 {
		p.Mime = mime[0]
	}
	if p.Template == "" {
		if p.Mime == "video" || p.Mime == "image" {
			p.SetTemplate("swiper")
		}
	}
	p.Files = append(p.Files, GlobMime(p.Path, p.Mime)...)
	return p
}

func (p *Page) GlobExt(ext ...string) *Page {
	p.glob = Extension
	if len(ext) > 0 {
		p.Mime = mime.TypeByExtension(ext[0])
		if strings.Contains(p.Mime, "video") {
			p.Mime = "video"
			p.SetTemplate("swiper")
		}
		if strings.Contains(p.Mime, "image") {
			p.Mime = "image"
			p.SetTemplate("swiper")
		}
	}
	p.Ext = ext
	p.Files = append(p.Files, GlobExt(p.Path, p.Ext...)...)
	return p
}

func (p *Page) GetChildren() *Page {
	entries := GetDirEntries(p.Path)
	for _, entry := range entries {
		fp := filepath.Join(p.Path, entry.Name())
		if entry.IsDir() {
			child := NewPage(fp)
			switch p.glob {
			case MimeType:
				child.GlobMime(p.Mime)
			case Extension:
				child.GlobExt(p.Ext...)
			}
			p.Children = append(p.Children, child.GetChildren())
		}
		switch name := entry.Name(); name {
		case "meta.toml":
			t, err := os.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}
			toml.Unmarshal(t, &p.Meta)
		}
	}
	return p
}

func (p Page) Title() string {
	return filepath.Base(p.Path)
}

func (p Page) HasChildren() bool {
	return len(p.Children) > 0
}

func (p *Page) Render() []byte {
	if p.Template == "swiper" && p.category == "" {
		p.Category.Html = GetCategory("swiper").Html
	}

	var buf bytes.Buffer
	err := Templates.ExecuteTemplate(&buf, "base", p)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (p *Page) SetTemplate(t string) *Page {
	p.Template = t
	return p
}

func (p *Page) Content() string {
	var (
		buf bytes.Buffer
		err error
	)

	switch p.Template {
	case "":
		fallthrough
	case "filetree":
		c := NewCollection(p.Path)
		c.GlobMime("").GetChildren()
		return c.Content()
	case "swiper":
		err = Templates.ExecuteTemplate(&buf, "swiper", p)
	default:
		t := template.Must(template.New("content").ParseFiles(p.Template))
		err = t.ExecuteTemplate(&buf, "content", p)
	}

	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func (p *Page) SetUrl(u string) {
	p.Url = u
}
