package static

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
)

type Index interface {
	Title() string
	Content() string
}

type Html struct {
	Scripts []string `toml:"scripts"`
	Css     []string `toml:"css"`
	Video   Video    `toml:"video"`
	Index
}

type Video struct {
	Muted    bool `toml:"muted"`
	Controls bool `toml:"controls"`
	Loop     bool `toml:"loop"`
	Autoplay bool `toml:"autoplay"`
}

func (html Html) RenderPage(p Index) []byte {
	html.Index = p

	var buf bytes.Buffer
	err := Templates.ExecuteTemplate(&buf, "base", html)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func (h Html) BaseCss() []string {
	var css []string
	for _, c := range h.Css {
		css = append(css, filepath.Base(c))
	}
	return css
}

func (h Html) BaseScripts() []string {
	var scripts []string
	for _, s := range h.Scripts {
		scripts = append(scripts, filepath.Base(s))
	}
	return scripts
}

func (c *Html) AddScripts(scripts ...string) {
	c.Scripts = append(c.Scripts, scripts...)
}

func (c *Html) AddCss(css ...string) {
	c.Css = append(c.Css, css...)
}

func (c Html) ReadScripts() []string {
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

func (c Html) ReadCss() []string {
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
