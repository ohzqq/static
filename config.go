package static

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"golang.org/x/exp/slices"
)

var (
	Opts    Config
	Default Config
)

type Color struct {
	Bg     string
	Fg     string
	White  string
	Black  string
	Grey   string
	Yellow string
	Red    string
	Pink   string
	Cyan   string
	Blue   string
	Green  string
	Purple string
}

type Config struct {
	Html
	Path       string
	Categories []string            `toml:"categories"`
	Color      Color               `toml:"color"`
	Category   map[string]Category `toml:"category"`
}

func ParseConfig(path string) (Config, error) {
	var (
		cfg  Config
		data []byte
		err  error
	)

	switch path {
	case "static/config.toml":
		data, err = Static.ReadFile(path)
		if err != nil {
			return cfg, err
		}
	default:
		data, err = os.ReadFile(path)
		if err != nil {
			return cfg, err
		}
	}

	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}
	cfg.Path = path

	dir := filepath.Dir(path)
	cfg.Scripts = AbsolutePaths(dir, cfg.Scripts...)
	cfg.Css = AbsolutePaths(dir, cfg.Css...)

	for name, col := range cfg.Category {
		if len(col.Scripts) > 0 {
			col.Scripts = AbsolutePaths(dir, col.Scripts...)
		}
		if len(col.Css) > 0 {
			col.Css = AbsolutePaths(dir, col.Css...)
		}
		if col.Template != "" {
			col.Template = AbsolutePaths(dir, col.Template)[0]
		}
		col.Css = append(cfg.Css, col.Css...)
		col.Scripts = append(cfg.Scripts, col.Scripts...)
		cfg.Category[name] = col
	}

	Opts = cfg

	return Opts, nil
}

func ParseDefault() Config {
	cfg, err := ParseConfig("static/config.toml")
	if err != nil {
		log.Fatal(err)
	}

	Default = cfg
	return cfg
}

func AbsolutePaths(root string, path ...string) []string {
	var paths []string
	for _, p := range path {
		switch filepath.IsAbs(p) {
		case true:
			paths = append(paths, p)
		case false:
			abs := filepath.Join(root, p)
			paths = append(paths, abs)
		}
	}
	return paths
}

func DefaultHtml() Html {
	html := Html{
		Css:     Default.Css,
		Scripts: Default.Scripts,
	}

	for _, css := range Opts.Css {
		b := filepath.Base(css)
		if !slices.Contains(html.BaseCss(), b) {
			html.AddCss(css)
		}
	}

	for _, script := range Opts.Scripts {
		if !slices.Contains(html.BaseScripts(), filepath.Base(script)) {
			html.AddScripts(script)
		}
	}

	return html
}

func Colors() Color {
	var c Color
	if Opts.Color != c {
		return Opts.Color
	}
	return Default.Color
}
