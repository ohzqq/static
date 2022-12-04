package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Path       string
	Color      Color                 `toml:"color"`
	Categories []string              `toml:"categories"`
	Scripts    []string              `toml:"scripts"`
	Css        []string              `toml:"css"`
	Collection map[string]Collection `toml:"collection"`
}

var Opts Config

func ParseConfig(path string) error {
	t, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(t, &Opts)
	if err != nil {
		return err
	}
	Opts.Path = path

	dir := filepath.Dir(path)
	Opts.Scripts = AbsolutePaths(dir, Opts.Scripts...)
	Opts.Css = AbsolutePaths(dir, Opts.Css...)

	for name, col := range Collections() {
		col.Scripts = AbsolutePaths(dir, col.Scripts...)
		col.Css = AbsolutePaths(dir, col.Css...)
		col.Template = AbsolutePaths(dir, col.Template)[0]
		col.Css = append(Opts.Css, col.Css...)
		//col.AddCss(Opts.Css...)
		col.AddScripts(Opts.Scripts...)
		Opts.Collection[name] = col
	}

	return nil
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
