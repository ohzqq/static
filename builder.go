package static

import (
	"io/fs"
	"log"
	"strings"
	"text/template"

	"github.com/ohzqq/fidi"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

type Site struct {
	Input string
}

func New(path string) *Site {
	return &Site{
		Input: path,
	}
}

func (b *Site) Build() {
	if b.Input == "" {
		log.Fatal("no input")
	}
	tree := fidi.NewTree(b.Input)

	page := NewPage(tree)
	page.Build()

	if recurse() {
		for _, child := range page.Children {
			child.Build()
		}
	}
}

func parsePageResources(kind string) []string {
	files := viper.GetStringSlice("global" + kind)

	if pro := viper.GetString("build.profile"); pro != "global" {
		if in := inheritedProfile(pro); in != "" {
			inh := viper.GetStringSlice(in + kind)
			files = append(files, inh...)
		}

		proF := viper.GetStringSlice(pro + kind)
		files = append(files, proF...)
	}

	return readScriptsAndStyles(files)
}

func getTemplate() *template.Template {
	tmpl := Templates.Lookup("filterableList")
	if hasMimes() {
		tmpl = Templates.Lookup("mediaAsset")
	}

	pro := viper.GetString("build.profile")
	if pro != "global" {
		if in := inheritedProfile(pro); in != "" {
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
		if in := inheritedProfile(pro); in != "" {
			inh := unmarshalHtml(in)
			maps.Copy(html, inh)
		}

		h := unmarshalHtml(pro)
		maps.Copy(html, h)
	}
	return html
}

func readScriptsAndStyles(files []string) []string {
	var assets []string
	for _, asset := range files {
		var f fs.FS
		if strings.HasPrefix(asset, "public") {
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
