package cmd

import (
	"log"
	"os"
	"path/filepath"
	"static"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	extension []string
	mimetype  string
	category  string
	cfg       static.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "static",
	Short: "A brief description of your application",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/static/config.toml)")
	rootCmd.PersistentFlags().StringSliceVarP(&extension, "ext", "e", []string{}, "glob by ext")
	rootCmd.PersistentFlags().StringVarP(&mimetype, "mimetype", "m", "", "glob by mimetype")
	rootCmd.PersistentFlags().StringVarP(&category, "category", "c", "", "config category")
	rootCmd.MarkFlagsMutuallyExclusive("ext", "mimetype")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".static" (without extension).
		path := filepath.Join(home, ".config", "static")
		viper.AddConfigPath(path)
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	static.ParseDefault()
	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		cfile := viper.ConfigFileUsed()
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatal(err)
		}
		cfg, err = static.ParseConfig(cfile)
		if err != nil {
			log.Fatal(err)
		}
	}
}
