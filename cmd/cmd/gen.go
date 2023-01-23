package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "(re)generate a site",
	Args:  cobra.ExactArgs(1),
	Run:   runGenCmd,
}

// genIndexCmd represents the gen command
var genIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "(re)generate a site",
	Args:  cobra.ExactArgs(1),
	Run:   runGenCmd,
}

// genCollectionCmd represents the gen command
var genCollectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "(re)generate a site",
	Args:  cobra.ExactArgs(1),
	Run:   runGenCmd,
}

// genAllCmd represents the gen command
var genAllCmd = &cobra.Command{
	Use:   "all",
	Short: "(re)generate a site",
	Args:  cobra.ExactArgs(1),
	Run:   runGenCmd,
}

func runGenCmd(cmd *cobra.Command, args []string) {
	input := args[0]
	parseFlags()
	viper.Set("build.regen", true)
	parseSubCommands(cmd)
	buildSite(input)
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.AddCommand(genIndexCmd)
	genCmd.AddCommand(genCollectionCmd)
	genCmd.AddCommand(genAllCmd)
}
