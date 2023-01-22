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

func parsePageResources(kind string) []string {
	files := viper.GetStringSlice("global" + kind)

	if pro := viper.GetString("build.profile"); pro != "global" {
		if in := InheritedProfile(pro); in != "" {
			inh := viper.GetStringSlice(in + kind)
			files = append(files, inh...)
		}

		proF := viper.GetStringSlice(pro + kind)
		files = append(files, proF...)
	}

	return ReadScriptsAndStyles(files)
}

func GetTemplate() *template.Template {
	var tmpl *template.Template

	pro := viper.GetString("build.profile")
	if pro != "global" {
		if in := InheritedProfile(pro); in != "" {
			pro = in
		}
		tmpl = Templates.Lookup(pro)
		if tmpl == nil {
			log.Fatalf("template %s not found\n", pro)
		}
	}

	return tmpl
}

func parseFilterKind(kind string) []string {
	pro := viper.GetString("build.profile")
	if pro != "global" {
		return viper.GetStringSlice(pro + kind)
	}
	return viper.GetStringSlice("global" + kind)
}

func getHtml() Html {
	html := unmarshalHtml("global")
	if pro := viper.GetString("build.profile"); pro != "global" {
		if in := InheritedProfile(pro); in != "" {
			inh := unmarshalHtml(in)
			maps.Copy(html, inh)
		}

		h := unmarshalHtml(pro)
		maps.Copy(html, h)
	}
	return html
}

func unmarshalHtml(pro string) Html {
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
