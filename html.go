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
	Video   `toml:"video"`
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

func (h Html) ReadScripts() []string {
	var scripts []string
	for _, script := range h.Scripts {
		scripts = append(scripts, ReadAsset(script))
	}
	return scripts
}

func (h Html) ReadCss() []string {
	var css []string
	for _, style := range h.Css {
		css = append(css, ReadAsset(style))
	}
	return css
}

func ReadAsset(name string) string {
	var (
		data []byte
		err  error
	)

	if filepath.IsAbs(name) {
		data, err = os.ReadFile(name)
	} else {
		data, err = Static.ReadFile(name)
	}

	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}
