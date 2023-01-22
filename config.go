package static

import (
	"embed"
	"io/fs"
	"log"
	"strings"
	"text/template"

	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

var (
	//go:embed static
	Public  embed.FS
	UserCfg fs.FS
	//defaultProfile Profile
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

type Html map[string]map[string]any

func Profiles() []string {
	cfg := viper.AllSettings()
	keys := maps.Keys(cfg)

	var profiles []string
	for _, k := range keys {
		if k != "color" && k != "build" {
			profiles = append(profiles, k)
		}
	}

	return profiles
}

func InheritedProfile(pro string) string {
	p := pro + ".inherit"
	if viper.IsSet(p) {
		return viper.GetString(p)
	}
	return ""
}

func parseCss(pro string) []string {
	var files []string

	if in := InheritedProfile(pro); in != "" {
		files = viper.GetStringSlice(in + ".css")
	}

	proF := viper.GetStringSlice(pro + ".css")
	files = append(files, proF...)

	return ReadScriptsAndStyles(files)
}

func parseScripts(pro string) []string {
	var files []string

	if in := InheritedProfile(pro); in != "" {
		files = viper.GetStringSlice(in + ".scripts")
	}

	proF := viper.GetStringSlice(pro + ".scripts")
	files = append(files, proF...)

	return ReadScriptsAndStyles(files)
}

func GetTemplate(pro string) *template.Template {
	if viper.IsSet("inherits." + pro) {
		pro = InheritedProfile(pro)
	}

	tmpl := Templates.Lookup(pro)
	if tmpl == nil {
		log.Fatalf("template %s not found\n", pro)
	}

	return tmpl
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
