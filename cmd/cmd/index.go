package cmd

import (
	"static"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "index of files in dir",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		parseFlags()
		site := static.New(input)
		viper.Set("build.index_only", true)
		site.Build()
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
