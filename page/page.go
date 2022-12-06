package page

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"static/config"
	"static/files"

	"github.com/BurntSushi/toml"
)

type Meta struct {
	Title string   `toml:"title"`
	Tags  []string `toml:"tags"`
}

type Page struct {
	Meta     Meta
	Type     string
	Url      string
	Path     string `toml:"path"`
	Files    []string
	Children []*Page
	Recurse  bool
	config.Collection
}

func NewPage(root string, collection config.Collection) *Page {
	page := Page{
		Path:       root,
		Collection: collection,
	}
	page.Files = append(page.Files, files.GlobMime(page.Path, page.Mime)...)

	return &page
}

func NewPageWithChildren(root string, collection config.Collection) *Page {
	page := NewPage(root, collection)
	page.Url = "./index.html"
	page.MakeIndexWithMime()
	return page
}

func (p *Page) MakeIndexWithMime() *Page {
	entries := files.GetDirEntries(p.Path)

	for _, e := range entries {
		fp := filepath.Join(p.Path, e.Name())
		if e.IsDir() {
			child := NewPage(fp, p.Collection)
			child.MakeIndexWithMime()
			p.Children = append(p.Children, child)
		}
		switch name := e.Name(); name {
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

	if p.Template != "" {
		t := template.Must(template.New("content").ParseFiles(p.Template))
		err := t.ExecuteTemplate(&buf, "content", p)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var err error
		switch p.Mime {
		case "video", "image":
			err = Templates.ExecuteTemplate(&buf, "swiper", p)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	return buf.String()
}

func (p *Page) SetUrl(u string) {
	p.Url = u
}
