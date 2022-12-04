package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Categories []string              `toml:"categories"`
	Scripts    []string              `toml:"scripts"`
	Css        []string              `toml:"css"`
	Collection map[string]Collection `toml:"collection"`
}

type Collection struct {
	Ext     []string `toml:"ext"`
	Scripts []string `toml:"scripts"`
	Css     []string `toml:"css"`
	Mime    string   `toml:"mime"`
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
	return nil
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

func GetCollection(col string) Collection {
	if c, ok := Opts.Collection[col]; ok {
		return c
	}
	return Collection{}
}
