package cmd

import (
	"log"
	"static"

	"github.com/spf13/cobra"
)

var (
	recurse  bool
	children bool
)

// pageCmd represents the page command
var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cat := static.GetCategory(args[0])
		dir := args[1]
		p := static.NewPage(dir).GlobMime(cat.Mime)
		if cmd.Flags().Changed("ext") {
			p = static.NewPage(dir).GlobExt(extension...).GetChildren()
		}
		err := static.Write(p.Path, cat.RenderPage(p))
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pageCmd)
}
