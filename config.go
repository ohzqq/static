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
	Public         embed.FS
	UserCfg        fs.FS
	defaultProfile Profile
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

type Html map[string]map[string]any

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

func ProfileInherits(pro string) bool {
	return viper.IsSet(pro + ".inherit")
}

func InheritedProfile(pro string) string {
	return viper.GetString(pro + ".inherit")
}

func GetCss(pro string) []string {
	var files []string

	if viper.IsSet(pro + ".inherit") {
		in := InheritedProfile(pro)
		files = viper.GetStringSlice(in + ".css")
	}

	proF := viper.GetStringSlice(pro + ".css")
	files = append(files, proF...)

	return ReadScriptsAndStyles(files)
}

func GetScripts(pro string) []string {
	var files []string

	if ProfileInherits(pro) {
		in := viper.GetString(pro + ".inherit")
		files = viper.GetStringSlice(in + ".scripts")
	}

	proF := viper.GetStringSlice(pro + ".scripts")
	files = append(files, proF...)

	return ReadScriptsAndStyles(files)
}

func GetHtml(pro string) Html {
	var html Html
	err := viper.UnmarshalKey(pro+".html", &html)
	if err != nil {
		log.Fatal(err)
	}
	return html
}

func ReadScriptsAndStyles(files []string) []string {
	var assets []string
	for _, asset := range files {
		var f fs.FS
		if strings.HasPrefix(asset, "static") {
			f = Public
		} else {
			f = UserCfg
		}
		d, err := fs.ReadFile(f, asset)
		if err != nil {
			log.Fatal(err)
		}
		assets = append(assets, string(d))
	}
	return assets
}
