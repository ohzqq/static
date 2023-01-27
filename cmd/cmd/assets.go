package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// assetsCmd represents the assets command
var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "export asset metadata",
	Args:  cobra.ExactArgs(1),
	Run:   runAssetsCmd,
}

// assetsGenCmd represents the gen command
var assetsGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "regen asset metadata",
	Args:  cobra.ExactArgs(1),
	Run:   runAssetsCmd,
}

// assetsCollectionCmd represents the gen command
var assetsCollectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "recursively export asset meta",
	Args:  cobra.ExactArgs(1),
	Run:   runAssetsCmd,
}

func runAssetsCmd(cmd *cobra.Command, args []string) {
	input := args[0]
	parseFlags()
	viper.Set("build.assets", true)
	viper.Set("build.no_thumbs", true)
	viper.Set("build.regen", true)
	parseSubCommands(cmd)
	buildSite(input)
}

func init() {
	rootCmd.AddCommand(assetsCmd)
	assetsCmd.AddCommand(assetsGenCmd)
	assetsCmd.AddCommand(assetsCollectionCmd)
	//assetsCmd.AddCommand(assetsgenCmd)
}
