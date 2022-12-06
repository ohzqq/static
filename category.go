package static

import (
	"bytes"
	"log"
	"os"
)

type Index interface {
	Title() string
	Content() string
}

type Category struct {
	Ext      []string `toml:"ext"`
	Scripts  []string `toml:"scripts"`
	Css      []string `toml:"css"`
	Mime     string   `toml:"mime"`
	Template string   `toml:"template"`
	Html     Html     `toml:"html"`
	Index
	//*Page
}

func (c Category) RenderPage(p Index) []byte {
	c.Index = p

	var buf bytes.Buffer
	err := Templates.ExecuteTemplate(&buf, "base", c)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (c Category) RecursiveWrite(pages ...*Page) error {
	for _, p := range pages {
		err := Write(p.Path, c.RenderPage(p))
		if err != nil {
			return err
		}

		if p.HasChildren() {
			err := c.RecursiveWrite(p.Children...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Category) AddScripts(scripts ...string) {
	c.Scripts = append(c.Scripts, scripts...)
}

func (c *Category) AddCss(css ...string) {
	c.Css = append(c.Css, css...)
}

func (c Category) ReadScripts() []string {
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

func (c Category) ReadCss() []string {
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
