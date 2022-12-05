package page

import (
	"bytes"
	"fmt"
	"html/template"
	"idxgen/config"
	"idxgen/files"
	"log"
	"os"
	"path/filepath"
)

type Meta struct {
	Title string   `toml:"title"`
	Tags  []string `toml:"tags"`
}

type Page struct {
	Meta     Meta
	Type     string
	Path     string `toml:"path"`
	Files    []string
	Children []Page
	config.Collection
}

func NewPage(root, collection string) Page {
	page := Page{
		Path:       root,
		Collection: config.GetCollection(collection),
	}
	page.Files = append(page.Files, files.GlobMime(page.Path, page.Mime)...)

	return page
}

func NewPageWithChildren(root, col string) Page {
	collection := config.GetCollection(col)
	page := MakeIndexWithMime(root, collection)
	page.Collection = collection
	page.Files = append(page.Files, files.GlobMime(page.Path, page.Mime)...)

	return page
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

func Write(pages ...Page) error {
	for _, page := range pages {
		out := filepath.Join(page.Path, "index.html")

		err := os.WriteFile(out, page.Render(), 0666)
		if err != nil {
			return fmt.Errorf("Rendering %s failed with error %s\n", out, err)
		}
		fmt.Printf("Rendered %s\n", out)
	}
	return nil
}

func RecursiveWrite(pages ...Page) error {
	for _, page := range pages {
		err := Write(page)
		if err != nil {
			return err
		}

		if page.HasChildren() {
			err := RecursiveWrite(page.Children...)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
		err := Templates.ExecuteTemplate(&buf, "body"+p.Mime, p)
		if err != nil {
			log.Fatal(err)
		}
	}

	return buf.String()
}
