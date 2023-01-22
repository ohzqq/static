package cmd

import (
	"fmt"

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
		fmt.Printf("regen %s\n", viper.GetBool("build.regen"))
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
