package cmd

import (
	"log"
	"static"

	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "render file tree for dir",
	Long:  "Command expects to find an index.html file in each sub-directory.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		p := static.NewCollection(dir)
		p.GlobMime("").GetChildren()
		html := static.DefaultHtml()
		err := static.Write(p.Path, html.RenderPage(p))
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
