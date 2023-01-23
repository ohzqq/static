package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build a static site",
	Args:  cobra.ExactArgs(1),
	Run:   runBuildCmd,
}

// buildCollectionCmd represents the collection command
var buildCollectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "build a collection",
	Args:  cobra.ExactArgs(1),
	Run:   runBuildCmd,
}

// buildAllCmd represents the collection command
var buildAllCmd = &cobra.Command{
	Use:   "all",
	Short: "build a collection",
	Args:  cobra.ExactArgs(1),
	Run:   runBuildCmd,
}

func runBuildCmd(cmd *cobra.Command, args []string) {
	input := args[0]
	parseFlags()
	viper.Set("build.index_only", false)
	parseSubCommands(cmd)
	buildSite(input)
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.AddCommand(buildCollectionCmd)
	buildCmd.AddCommand(buildAllCmd)
}
