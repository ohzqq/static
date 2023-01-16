package static

import (
	"embed"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

var (
	Opts    Config
	Default Config
)

var (
	//go:embed static/*
	Public embed.FS
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
	fmt.Printf("default %+V\n", defaultProfile)
	p := viper.Sub(pro)

	var profile Profile
	err := p.Unmarshal(&profile)
	if err != nil {
		log.Fatal(err)
	}
	return profile
}

func MergeProfiles(pro1, pro2 Profile) Profile {
	var pro Profile
}

func SetDefaultProfile() {
	p := viper.Sub("global")

	err := p.Unmarshal(&defaultProfile)
	if err != nil {
		log.Fatal(err)
	}
}
