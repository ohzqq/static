package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		Toot()
	},
}

func Toot(s ...string) {
	fmt.Printf("%+V\n", s)
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
