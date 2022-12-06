package static

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var (
	Opts    Config
	Default Config
)

type Config struct {
	Path       string
	Color      Color               `toml:"color"`
	Categories []string            `toml:"categories"`
	Scripts    []string            `toml:"scripts"`
	Css        []string            `toml:"css"`
	Category   map[string]Category `toml:"collection"`
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

func Categories() []string {
	return Opts.Categories
}

func Scripts() []string {
	return Opts.Scripts
}

func Css() []string {
	return Opts.Css
}

func Collections() map[string]Category {
	col := Default.Category
	for n, c := range Opts.Category {
		col[n] = c
	}
	return col
}

func Colors() Color {
	var c Color
	if Opts.Color != c {
		return Opts.Color
	}
	return Default.Color
}

func GetCollection(collection string) Category {
	if c, ok := Opts.Category[collection]; ok {
		return c
	}

	if c, ok := Default.Category[collection]; ok {
		return c
	}
	return Category{}
}
