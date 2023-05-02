package cmd

import (
	"encoding/json"
	"os"
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

		var idx []*static.Media
		err = json.Unmarshal(f, &idx)
		if err != nil {
			panic(err)
		}

		for _, m := range idx {
			var thumb []byte
			switch {
			case m.Video != "":
				thumb = static.VideoThumb(strings.TrimPrefix(m.Video, "/"))
			case m.Img != "":
				thumb = static.ImageThumb(strings.TrimPrefix(m.Img, "/"))
			}
			m.Thumbnail = static.ThumbToBase64(thumb)
		}

		gal, err := json.MarshalIndent(idx, "", "  ")
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
