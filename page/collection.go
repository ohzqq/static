package page

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"static/category"
	"static/config"
	"static/files"
)

type Collection struct {
	*Page
	Root     string
	Filetree string
	Ext      []string    `toml:"ext"`
	Scripts  []string    `toml:"scripts"`
	Css      []string    `toml:"css"`
	Mime     string      `toml:"mime"`
	Template string      `toml:"template"`
	Html     config.Html `toml:"html"`
}

func NewCollection(root string, collection category.Category) Collection {
	col := Collection{
		Root: root,
		Page: NewPage(root, collection),
	}
	col.Url = "./index.html"
	for _, p := range col.Children {
		RelativeUrls(root, p)
	}
	col.Filetree = col.Tree()
	return col
}

func (c Collection) Tree() string {
	return c.Content()
}

func (c Collection) Content() string {
	var buf bytes.Buffer
	err := Templates.ExecuteTemplate(&buf, "filetree", c)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (c Collection) Render() []byte {
	var buf bytes.Buffer
	err := Templates.ExecuteTemplate(&buf, "base", c)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func RelativeUrls(root string, pages ...*Page) []*Page {
	for _, page := range pages {
		rel, err := filepath.Rel(root, page.Path)
		if err != nil {
			log.Fatal(err)
		}
		u := filepath.Join(rel, "index.html")
		page.SetUrl("./" + u)

		if page.HasChildren() {
			page.Children = RelativeUrls(root, page.Children...)
		}
	}
	return pages
}

func NewCollectionWithExt(root string, ext ...string) *Collection {
	page := Collection{Root: root}
	page.Page = MakeIndexWithExt(root, ext...)
	page.Files = append(page.Files, files.GlobExt(page.Path, ext...)...)

	return &page
}

func (c *Collection) AddScripts(scripts ...string) {
	c.Scripts = append(c.Scripts, scripts...)
}

func (c *Collection) AddCss(css ...string) {
	c.Css = append(c.Css, css...)
}

func (c Collection) ReadScripts() []string {
	var scripts []string
	for _, s := range c.Scripts {
		t, err := os.ReadFile(s)
		if err != nil {
			log.Fatal(err)
		}
		scripts = append(scripts, string(t))
	}
	return scripts
}

func (c Collection) ReadCss() []string {
	var css []string
	for _, s := range c.Css {
		t, err := os.ReadFile(s)
		if err != nil {
			log.Fatal(err)
		}
		css = append(css, string(t))
	}
	return css
}
