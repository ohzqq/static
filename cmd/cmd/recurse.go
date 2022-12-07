package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// recurseCmd represents the recurse command
var recurseCmd = &cobra.Command{
	Use:   "recurse",
	Short: "recursively generate static, self-contained html pages",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		p := MakePage(dir, cmd)

		p.GetChildren()
		err := p.Category.RecursiveWrite(p)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(recurseCmd)
}
