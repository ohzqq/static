package cmd

import (
	"fmt"
	"static"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		d := args[0]
		col := static.NewCollection(d)
		fmt.Printf("%s\n", col.Index.Rel())
		//node := col.Nodes[2]
		//node, _ := col.GetNode(2)
		//fmt.Printf("cur %d\n", node.Rel())

		for _, page := range col.Pages() {
			page.FilterByMime("image")
			fmt.Printf("%d: %s\n", page.Info().Depth, page.Info().Rel())
			if page.HasIndex {
				fmt.Printf("index %s\n", page.Index.Rel())
			}
			fmt.Printf("assets %+V\n", page.Assets)
			//page := static.NewPage(node)
			for _, child := range page.Parents() {
				fmt.Printf("parent %+V\n", child.Info().Rel())
			}
		}

		//for _, page := range static.GetParentsByDepth(col.Tree, 2) {
		//  fmt.Printf("depth %s: %s\n", page.Depth, page.Rel())
		//}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
