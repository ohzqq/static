package cmd

import (
	"static"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		parseFlags()
		site := static.New(input)
		viper.Set("build.index_only", false)
		site.Build()
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

}
