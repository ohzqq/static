package cmd

import (
	"static"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func parseFlags() {
	switch {
	case rootCmd.PersistentFlags().Changed("all"), rootCmd.PersistentFlags().Changed("ext"), rootCmd.PersistentFlags().Changed("mime"), rootCmd.PersistentFlags().Changed("profile"):
		viper.Set("build.index_only", false)
	}
}

func parseSubCommands(cmd *cobra.Command) {
	switch cmd.Name() {
	case "all":
		viper.Set("build.all", true)
		viper.Set("build.index_only", true)
	case "collection":
		viper.Set("build.is_collection", true)
	case "index":
		viper.Set("build.index_only", true)
	case "gen":
		viper.Set("build.regen", true)
	case "assets":
		viper.Set("build.assets", true)
		viper.Set("build.no_thumbs", true)
	}
}

func buildSite(input string) {
	site := static.New(input)
	site.Build()
}
