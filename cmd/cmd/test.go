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
		fmt.Printf("arg %s\n", d)
		//dir, _ := fidi.NewDir(d, d)
		//p := static.NewPage(dir)
		static.GetTemplates()

		//collection(d)

	},
}

func collection(d string) {
	col := static.NewCollection(d, "gifv")
	fmt.Printf("collection %s\n", col.Title)
	fmt.Printf("collection %s\n", len(col.Css))
	//node := col.Nodes[2]
	//node, _ := col.GetNode(2)
	//pro := static.GetProfile("gifv")
	//fmt.Printf("cfg %v\n", pro)

	for _, page := range col.Children {
		//page.FilterByExt(".jpg", ".png", ".avif")
		fmt.Printf("%d: %s\n", page.Info().Depth, page.Info().Rel())
		if page.HasIndex {
			fmt.Printf("url %+V\n", page.RelUrl())
		}
		//fmt.Printf("breadcrumbs %+V\n", page.Breadcrumbs())
		//fmt.Printf("nav %s\n", page.Nav())
		//page := static.NewPage(node)
		for _, child := range page.Assets {
			fmt.Printf("asset %+V\n", child.Attributes)
			fmt.Printf("html %+V\n", child.Render())
		}
	}

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
