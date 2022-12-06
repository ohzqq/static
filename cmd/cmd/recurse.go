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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cat := static.GetCategory(category)
		dir := args[0]
		p := static.NewPage(dir)
		if cat.Mime != "" || cmd.Flags().Changed("mimetype") {
			p.GlobMime(cat.Mime)
		} else if len(cat.Ext) > 0 || cmd.Flags().Changed("ext") {
			p.GlobExt(extension...)
		}
		p.GetChildren()
		err := cat.RecursiveWrite(p)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(recurseCmd)
}
