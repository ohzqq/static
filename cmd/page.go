package cmd

import (
	"fmt"
	"idxgen/page"
	"log"
	"os"
	"path/filepath"

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
		collection := args[0]
		dir := args[1]
		p := page.NewPageWithChildren(dir, collection)
		out := filepath.Join(dir, "index.html")
		err := os.WriteFile(out, p.Render(), 0600)
		if err != nil {
			log.Fatalf("Rendering %s failed with error %s\n", out, err)
		}
		fmt.Printf("Rendered %s\n", out)
		//fmt.Printf("%+V\n", config.GetCollection(collection).Html.Video)
	},
}

func init() {
	rootCmd.AddCommand(pageCmd)
	pageCmd.Flags().BoolVarP(&recurse, "recurse", "r", false, "recursively generate index files")
	pageCmd.Flags().BoolVarP(&children, "children", "c", false, "list child directories")
}
