package static

import (
	"log"

	"github.com/ohzqq/fidi"
	"github.com/spf13/viper"
)

type Builder struct {
	Input string
}

func New(path string) *Builder {
	return &Builder{
		Input: path,
	}
}

func (b *Builder) Build() {
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
