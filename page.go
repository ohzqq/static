package static

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

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
	Mime     string
	Ext      []string
	Url      string
	Path     string `toml:"path"`
	Files    []string
	Children []*Page
	Recurse  bool
}

func NewPage(root string) *Page {
	page := Page{
		Path: root,
	}
	return &page
}

func (p *Page) GlobMime(mime ...string) *Page {
	p.glob = MimeType
	if len(mime) > 0 {
		p.Mime = mime[0]
	}
	p.Files = append(p.Files, GlobMime(p.Path, p.Mime)...)
	return p
}

func (p *Page) GlobExt(ext ...string) *Page {
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

func (p Page) Render() []byte {
	var buf bytes.Buffer
	err := Templates.ExecuteTemplate(&buf, "base", p)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (p Page) Content() string {
	var buf bytes.Buffer

	//if p.Template != "" {
	//t := template.Must(template.New("content").ParseFiles(p.Template))
	//err := t.ExecuteTemplate(&buf, "content", p)
	//if err != nil {
	//log.Fatal(err)
	//}
	//} else {
	var err error
	switch p.Mime {
	case "video", "image":
		err = Templates.ExecuteTemplate(&buf, "swiper", p)
	}
	if err != nil {
		log.Fatal(err)
	}
	//}

	return buf.String()
}

func (p *Page) SetUrl(u string) {
	p.Url = u
}
