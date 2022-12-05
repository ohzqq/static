package cmd

import (
	"fmt"
	"idxgen/page"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
