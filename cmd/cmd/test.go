package cmd

import (
	"fmt"
	"static"

	"github.com/ohzqq/fidi"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		d := args[0]
		//fmt.Printf("arg %s\n", d)
		//tn := static.ExtractThumbFromVideo(fidi.NewFile(d))
		//println(tn)

		//page(d)
		//collection(d)
		p := static.New(d)
		p.Build()
		//p.Profile = "swiper"
		//p.Regen = true
		//p.ListAll = true
		//p.Collection().Build()

	},
}

func page(d string) {
	dir, _ := fidi.NewDir(d, d)
	p := static.NewPage(dir)
	p.Build()
	fmt.Printf("html %+V\n", p.Info().Path())

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
