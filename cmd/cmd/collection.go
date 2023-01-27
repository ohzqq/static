package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// collectionCmd represents the collection command
var collectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "recursively build a site",
	Args:  cobra.ExactArgs(1),
	Run:   runCollectionCmd,
}

// colIndexCmd represents the gen command
var colIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "recursively collect indices",
	Args:  cobra.ExactArgs(1),
	Run:   runCollectionCmd,
}

// colAssetsCmd represents the gen command
var colAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "recursively collect assets",
	Args:  cobra.ExactArgs(1),
	Run:   runCollectionCmd,
}

// colGenCmd represents the gen command
var colGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "recursively generate pages",
	Args:  cobra.ExactArgs(1),
	Run:   runCollectionCmd,
}

func runCollectionCmd(cmd *cobra.Command, args []string) {
	input := args[0]
	parseFlags()
	viper.Set("build.is_collection", true)
	parseSubCommands(cmd)
	buildSite(input)
}

func init() {
	rootCmd.AddCommand(collectionCmd)
	collectionCmd.AddCommand(colIndexCmd)
	collectionCmd.AddCommand(colAssetsCmd)
	collectionCmd.AddCommand(colGenCmd)
}
