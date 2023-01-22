package cmd

import (
	"fmt"
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
		fmt.Printf("root dir %s\n", input)
		parseFlags()
		p := static.New(input)
		p.Page().Build()
	},
}

func parseFlags() {
	switch {
	case rootCmd.PersistentFlags().Changed("all"), rootCmd.PersistentFlags().Changed("ext"), rootCmd.PersistentFlags().Changed("mime"):
		viper.Set("build.index_only", false)
	}
}

func init() {
	rootCmd.AddCommand(buildCmd)

}
