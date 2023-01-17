package static

import (
	"embed"
	"io/fs"
	"log"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

var (
	//go:embed static
	Public  embed.FS
	UserCfg fs.FS
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

type Profile struct {
	Css     []string
	Scripts []string
	Html    map[string]map[string]any
}

func Profiles() []string {
	cfg := viper.AllSettings()
	keys := maps.Keys(cfg)

	var profiles []string
	for _, k := range keys {
		if k != "color" {
			profiles = append(profiles, k)
		}
	}

	return profiles
}

var defaultProfile Profile

func GetProfile(pro string) Profile {
	p := viper.Sub(pro)

	var profile Profile
	err := p.Unmarshal(&profile)
	if err != nil {
		log.Fatal(err)
	}

	merged := MergeProfiles(defaultProfile, profile)

	return merged
}

func (p Profile) ReadCss() []string {
	var assets []string
	for _, css := range p.Css {
		var f fs.FS
		if strings.HasPrefix(css, "static") {
			f = Public
		} else {
			f = UserCfg
		}
		d, err := fs.ReadFile(f, css)
		if err != nil {
			log.Fatal(err)
		}
		assets = append(assets, string(d))
	}
	return assets
}

func (p Profile) ReadScripts() []string {
	var assets []string
	for _, script := range p.Scripts {
		var f fs.FS
		if strings.HasPrefix(script, "static") {
			f = Public
		} else {
			f = UserCfg
		}
		d, err := fs.ReadFile(f, script)
		if err != nil {
			log.Fatal(err)
		}
		assets = append(assets, string(d))
	}
	return assets
}

func MergeProfiles(pro1, pro2 Profile) Profile {
	pro := Profile{
		Css:     pro1.Css,
		Scripts: pro1.Scripts,
		Html:    pro1.Html,
	}
	pro.Css = append(pro.Css, pro2.Css...)
	pro.Scripts = append(pro.Scripts, pro2.Scripts...)
	maps.Copy(pro.Html, pro2.Html)
	return pro
}

func SetDefaultProfile() {
	p := viper.Sub("global")

	err := p.Unmarshal(&defaultProfile)
	if err != nil {
		log.Fatal(err)
	}
}
