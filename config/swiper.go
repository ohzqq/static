package config

type Swiper struct {
	Video `toml:"video"`
}

type Video struct {
	Muted    bool `toml:"muted"`
	Controls bool `toml:"controls"`
	Loop     bool `toml:"loop"`
	Autoplay bool `toml:"autoplay"`
}
