package cmd

import (
	"log"
	"static"

	"github.com/spf13/cobra"
)

// pageCmd represents the page command
var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "generate a static, self-contained html page",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cat := static.GetCategory(category)
		dir := args[0]

		p := static.NewPage(dir)

		if cmd.Flags().Changed("category") {
			p.SetCategory(category)
		}

		if len(cat.Ext) > 0 || cmd.Flags().Changed("ext") {
			p.GlobExt(extension...)
		} else {
			var m string
			if cat.Mime != "" {
				m = cat.Mime
			}
			if cmd.Flags().Changed("mimetype") {
				m = mimetype
			}
			p.GlobMime(m)
		}

		err := static.Write(p.Path, p.Render())
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pageCmd)
}
