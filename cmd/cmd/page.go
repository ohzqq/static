package cmd

import (
	"log"
	"static"

	"github.com/spf13/cobra"
)

var (
	recurse   bool
	children  bool
	extension []string
)

// pageCmd represents the page command
var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		collection := static.GetCollection(args[0])
		dir := args[1]

		switch recurse {
		case true:
			p := static.New(dir).GlobMime(collection.Mime).GetChildren()
			if cmd.Flags().Changed("ext") {
				p = static.New(dir).GlobExt(extension...).GetChildren()
			}
			err := collection.RecursiveWrite(p)
			if err != nil {
				log.Fatal(err)
			}
		case false:
			p := static.New(dir).GlobMime(collection.Mime)
			err := static.Write(p.Path, collection.RenderPage(p))
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pageCmd)
	pageCmd.Flags().BoolVarP(&recurse, "recurse", "r", false, "recursively generate index files")
	pageCmd.Flags().BoolVarP(&children, "children", "c", false, "list child directories")
	pageCmd.Flags().StringSliceVarP(&extension, "ext", "e", []string{}, "list child directories")
}
