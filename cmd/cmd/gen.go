package cmd

import (
	"log"
	"static"

	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	//Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cat := static.GetCollection(args[0])
		dir := args[1]
		p := static.NewCollection(dir)
		//files := page.Batch(p.Files)
		//for _, pp := range p.Children {
		err := static.Write(p.Path, cat.RenderPage(p.Page))
		if err != nil {
			log.Fatal(err)
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
