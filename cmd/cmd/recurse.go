package cmd

import (
	"log"
	"static"

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

func MakePage(dir string, cmd *cobra.Command) *static.Page {
	cat := static.GetCategory(category)
	p := static.NewPage(dir)

	if cmd.Flags().Changed("category") {
		p.SetCategory(category)
	}

	if len(cat.Ext) > 0 || cmd.Flags().Changed("ext") {
		p.GlobExt(extension...)
	} else {
		var m string
		if cat.Mime != "" {
			m = cat.Mime
		}
		if cmd.Flags().Changed("mimetype") {
			m = mimetype
		}
		p.GlobMime(m)
	}

	return p
}
func init() {
	rootCmd.AddCommand(recurseCmd)
}
