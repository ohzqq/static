package cmd

import (
	"log"
	"static"

	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		p := static.NewCollection(dir)
		p.GlobMime("").GetChildren()
		//println(p.Content())
		err := static.Write(p.Path, static.DefaultHtml().RenderPage(p))
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
