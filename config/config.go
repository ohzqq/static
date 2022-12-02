package config

type Config struct {
	Categories []string              `toml:"categories"`
	Collection map[string]Collection `toml:"collection"`
}

type Collection struct {
	Ext []string `toml:"ext"`
}
