package cmd

import (
	"static"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// collectionCmd represents the collection command
var collectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run:   runCollectionCmd,
}

func runCollectionCmd(cmd *cobra.Command, args []string) {
	input := args[0]
	parseFlags()

	viper.Set("build.is_collection", true)

	site := static.New(input)
	site.Build()
}

func init() {
	rootCmd.AddCommand(collectionCmd)
}
