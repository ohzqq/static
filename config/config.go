package config

import (
	"idx"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var (
	Opts Config
)

type Config struct {
	Path       string
	Color      Color                 `toml:"color"`
	Categories []string              `toml:"categories"`
	Scripts    []string              `toml:"scripts"`
	Css        []string              `toml:"css"`
	Collection map[string]Collection `toml:"collection"`
}

func ParseConfig(path string) (Config, error) {
	var (
		cfg  Config
		data []byte
		err  error
	)

	switch path {
	case "default":
		data, err = idx.Static.ReadFile("static/config.toml")
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

	for name, col := range Collections() {
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
		cfg.Collection[name] = col
	}

	Opts = cfg

	return Opts, nil
}

func Default() Config {
	cfg, err := ParseConfig("default")
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func AbsolutePaths(root string, path ...string) []string {
	var paths []string
	for _, p := range path {
		switch filepath.IsAbs(p) {
		case true:
			paths = append(paths, p)
		case false:
			paths = append(paths, filepath.Join(root, p))
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

func Collections() map[string]Collection {
	return Opts.Collection
}

func Colors() Color {
	return Opts.Color
}

func GetCollection(col string) Collection {
	if c, ok := Opts.Collection[col]; ok {
		return c
	}
	return Collection{}
}
