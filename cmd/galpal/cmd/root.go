package cmd

import (
	"encoding/json"
	"mime"
	"os"
	"path/filepath"
	"static"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "galpal",
	Short: "generate gal json",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		files, err := os.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		var media []static.Media
		for _, file := range files {
			if !file.IsDir() {
				var m static.Media
				var thumb []byte
				name := filepath.Join(dir, file.Name())
				ext := filepath.Ext(name)
				mt := mime.TypeByExtension(ext)

				switch {
				case strings.Contains(mt, "video"):
					thumb = static.VideoThumb(name)
					m.Video = name
				case strings.Contains(mt, "image"):
					thumb = static.ImageThumb(name)
					m.Img = name
				}
				m.Thumbnail = static.ThumbToBase64(thumb)

				media = append(media, m)
			}
		}

		gal, err := json.MarshalIndent(media, "", "  ")
		if err != nil {
			panic(err)
		}

		idx, err := os.Create(filepath.Join(dir, "index.json"))
		if err != nil {
			panic(err)
		}
		defer idx.Close()

		_, err = idx.Write(gal)
		if err != nil {
			panic(err)
		}
	},
}

func checkForIndex(dir string) (string, bool) {
	idx, _ := filepath.Glob(filepath.Join(dir, "index.json"))

	if len(idx) > 0 {
		return idx[0], true
	}
	return "", false
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.galpal.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
