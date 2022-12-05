package cmd

import (
	"fmt"
	"idxgen/page"
	"path/filepath"

	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		//var err error

		for _, cat := range cfg.Categories {
			col := cfg.Collection[cat]
			path := filepath.Join(dir, cat)
			idx := page.NewCollectionWithExt(path, col.Ext...)
			idx.Collection = col
			idx.Type = cat

			fmt.Printf("%v idx %+V\n", idx.Title(), idx.Files)
			for _, c := range idx.Children {
				fmt.Printf("%v idx.Children %+V\n", c.Title(), c.Files)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
