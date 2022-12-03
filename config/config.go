package config

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
}
