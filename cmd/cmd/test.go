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
		fmt.Printf("%s\n", col.Path())
		node, _ := col.GetNode(2)
		fmt.Printf("cur %d\n", node.Rel())

		nodes := static.GetChildrenByDepth(col.Tree, 2)
		fmt.Printf("nodes %+V\n", len(nodes))

		for _, page := range nodes {
			//page := static.NewPage(node)
			fmt.Printf("depth %s: %s\n", page.Depth, page.Rel())
			//for _, child := range page.Children {
			//fmt.Printf("%+V\n", child.Rel())
			//}
		}

		for _, page := range static.GetParentsByDepth(col.Tree, 2) {
			fmt.Printf("depth %s: %s\n", page.Depth, page.Rel())
		}
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
