package cmd

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"static"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create new gal",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		var media []static.Media
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			var m static.Media
			var thumb []byte
			mt := static.MimeType(path)

			if mt.IsMedia() {
				switch {
				case mt.IsVideo():
					if cmd.Flags().Changed("thumbs") {
						thumb = static.VideoThumb(path)
						m.Thumbnail = static.ThumbToBase64(thumb)
					}
					m.Video = filepath.Join("/", path)
				case mt.IsImage():
					if cmd.Flags().Changed("thumbs") {
						thumb = static.ImageThumb(path)
						m.Thumbnail = static.ThumbToBase64(thumb)
					}
					m.Img = filepath.Join("/", path)
				}

				media = append(media, m)
			}
			return nil
		})

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

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolP("thumbs", "t", false, "create thumbs")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
