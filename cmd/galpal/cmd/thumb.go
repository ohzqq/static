package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"static"
	"strings"

	"github.com/spf13/cobra"
)

// thumbCmd represents the thumb command
var thumbCmd = &cobra.Command{
	Use:   "thumb",
	Short: "add thumbs to index",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.ReadFile(args[0])
		if err != nil {
			panic(err)
		}

		var gallery static.Gallery
		err = json.Unmarshal(f, &gallery)
		if err != nil {
			panic(err)
		}

		for _, m := range gallery.Media {
			var thumb []byte
			switch {
			case m.Video != "":
				name := strings.TrimPrefix(m.Video, "/")
				if args[0] == "index.json" {
					name = filepath.Base(name)
				}
				thumb = static.VideoThumb(name)
			case m.Img != "":
				name := strings.TrimPrefix(m.Img, "/")
				if args[0] == "index.json" {
					name = filepath.Base(name)
				}
				thumb = static.ImageThumb(name)
			}
			m.Thumbnail = static.ThumbToBase64(thumb)
		}

		gal, err := json.MarshalIndent(gallery, "", "  ")
		if err != nil {
			panic(err)
		}

		i, err := os.Create(args[0])
		if err != nil {
			panic(err)
		}
		defer i.Close()

		_, err = i.Write(gal)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(thumbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// thumbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// thumbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
