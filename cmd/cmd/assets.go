package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// assetsCmd represents the assets command
var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run:   runAssetsCmd,
}

func runAssetsCmd(cmd *cobra.Command, args []string) {
	input := args[0]
	parseFlags()
	viper.Set("build.assets", true)
	buildSite(input)
}

func init() {
	rootCmd.AddCommand(assetsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// assetsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// assetsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
