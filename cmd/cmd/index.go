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
	Run:   runIndexCmd,
}

// indexAllCmd represents the collection command
var indexAllCmd = &cobra.Command{
	Use:   "all",
	Short: "build a collection",
	Args:  cobra.ExactArgs(1),
	Run:   runIndexCmd,
}

func runIndexCmd(cmd *cobra.Command, args []string) {
	input := args[0]
	parseFlags()
	viper.Set("build.index_only", true)

	switch cmd.Name() {
	case "all":
		viper.Set("build.all", true)
	case "collection":
		viper.Set("build.is_collection", true)
	}

	site := static.New(input)
	site.Build()
}

func init() {
	rootCmd.AddCommand(indexCmd)
	buildCmd.AddCommand(indexCmd)
}
