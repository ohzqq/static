package static

import "embed"

var (
	Opts    Config
	Default Config
)

var (
	//go:embed static/*
	Static embed.FS
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
	//Html
	Path       string
	Categories []string `toml:"categories"`
	Color      Color    `toml:"color"`
	//Category   map[string]Category `toml:"category"`
}
