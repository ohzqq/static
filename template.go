package static

import (
	"fmt"
	"log"
	"path/filepath"
	"text/template"

	"github.com/spf13/viper"
)

var (
	TmplFuncs = template.FuncMap{
		"Batch":      AssetBatch,
		"Dir":        filepath.Dir,
		"Thumb":      Thumbnail,
		"VideoThumb": ExtractThumbFromVideo,
	}
)

var Templates = template.New("").Funcs(TmplFuncs)

func InitTemplates() []string {
	var (
		def  []string
		user []string
		err  error
	)

	for _, pro := range Profiles() {
		switch pro {
		case "swiper", "global":
			t := fmt.Sprintf("static/%s/*tmpl", pro)
			def = append(def, t)
		default:
			if !viper.IsSet(pro + ".inherit") {
				t := fmt.Sprintf("%s/*tmpl", pro)
				user = append(user, t)
			}
		}
	}

	Templates, err = Templates.ParseFS(Public, def...)
	if err != nil {
		log.Fatal(err)
	}

	Templates, err = Templates.ParseFS(UserCfg, user...)
	if err != nil {
		log.Fatal(err)
	}

	return def
}

func AssetBatch(og []Asset) [][]Asset {
	var (
		batch = 4
		files [][]Asset
	)
	for i := 0; i < len(og); i += batch {
		j := i + batch
		if j > len(og) {
			j = len(og)
		}

		files = append(files, og[i:j])
	}
	return files
}

func Batch(og []string) [][]string {
	var (
		batch = 4
		files [][]string
	)
	for i := 0; i < len(og); i += batch {
		j := i + batch
		if j > len(og) {
			j = len(og)
		}

		files = append(files, og[i:j])
	}
	return files
}
