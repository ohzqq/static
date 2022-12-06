package static

type PageElements struct {
	Html Html `toml:"html"`
}

type Html struct {
	Scripts []string `toml:"scripts"`
	Css     []string `toml:"css"`
	Video   Video    `toml:"video"`
}

type Video struct {
	Muted    bool `toml:"muted"`
	Controls bool `toml:"controls"`
	Loop     bool `toml:"loop"`
	Autoplay bool `toml:"autoplay"`
}
