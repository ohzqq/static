package cmd

import (
	"fmt"
	"log"
	"static"

	"github.com/spf13/cobra"
)

// pageCmd represents the page command
var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "generate a static, self-contained html page",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		p := MakePage(dir, cmd)
		fmt.Printf("%+V\n", static.Opts.Html)

		err := static.Write(p.Path, p.Render())
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pageCmd)
}
