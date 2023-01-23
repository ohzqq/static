package static

import (
	"embed"
	"io/fs"
	"log"

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

func inheritedProfile(pro string) string {
	p := pro + ".inherit"
	if viper.IsSet(p) {
		return viper.GetString(p)
	}
	return ""
}

func unmarshalHtml(pro string) Html {
	var html Html
	err := viper.UnmarshalKey(pro+".html", &html)
	if err != nil {
		log.Fatal(err)
	}
	return html
}

func listAll() bool {
	switch {
	case viper.GetBool("build.all"):
		return true
	case indexOnly():
		if recurse() {
			return false
		}
		return true
	default:
		return false
	}
}

func indexOnly() bool {
	return viper.GetBool("build.index_only")
}

func noThumbs() bool {
	return viper.GetBool("build.no_thumbs")
}

func recurse() bool {
	return viper.GetBool("build.is_collection")
}

func regen() bool {
	return viper.GetBool("build.regen")
}

func hasMimes() bool {
	return len(mimes()) > 0
}

func mimes() []string {
	return viper.GetStringSlice("build.mimes")
}

func hasExts() bool {
	return len(exts()) > 0
}

func exts() []string {
	return viper.GetStringSlice("build.exts")
}
