package cmd

import (
	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "index of files in dir",
	Args:  cobra.ExactArgs(1),
	Run:   runIndexCmd,
}

// indexCollectionCmd represents the collection command
var indexCollectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "build a collection",
	Args:  cobra.ExactArgs(1),
	Run:   runIndexCmd,
}

func runIndexCmd(cmd *cobra.Command, args []string) {
	input := args[0]
	parseFlags()
	parseSubCommands(cmd)
	buildSite(input)
}

func init() {
	rootCmd.AddCommand(indexCmd)
	indexCmd.AddCommand(indexCollectionCmd)
}
