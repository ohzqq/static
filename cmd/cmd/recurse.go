package cmd

import (
	"log"
	"static"

	"github.com/spf13/cobra"
)

// recurseCmd represents the recurse command
var recurseCmd = &cobra.Command{
	Use:   "recurse",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cat := static.GetCategory(args[0])
		dir := args[1]
		p := static.NewPage(dir).GlobMime(cat.Mime).GetChildren()
		if cmd.Flags().Changed("ext") {
			p = static.NewPage(dir).GlobExt(extension...).GetChildren()
		}
		err := cat.RecursiveWrite(p)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(recurseCmd)
}
