package config

import (
	"log"
	"os"
)

type Collection struct {
	Ext      []string `toml:"ext"`
	Scripts  []string `toml:"scripts"`
	Css      []string `toml:"css"`
	Mime     string   `toml:"mime"`
	Template string   `toml:"template"`
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
