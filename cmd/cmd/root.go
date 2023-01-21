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
	cfgFile    string
	extension  []string
	mimetype   string
	profile    string
	regenerate bool
	generate   bool
	cfg        static.Config
	Builder    *static.Builder
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

	rootCmd.PersistentFlags().BoolVarP(&regenerate, "regenerate", "G", false, "regenerate index files")
	rootCmd.PersistentFlags().BoolVarP(&generate, "generate", "g", false, "generate index files")
	rootCmd.PersistentFlags().StringSliceVarP(&extension, "ext", "e", []string{}, "glob by ext")
	rootCmd.PersistentFlags().StringVarP(&mimetype, "mimetype", "m", "", "glob by mimetype")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "config category")
	rootCmd.MarkFlagsMutuallyExclusive("ext", "mimetype")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault(
		"global.css",
		[]string{
			"static/css/normalize.css",
			"static/css/milligram.css",
			"static/css/base.css",
		},
	)

	viper.SetDefault(
		"global.scripts",
		[]string{
			"static/js/alpine.js",
		},
	)

	viper.SetDefault(
		"global.html.video",
		map[string]any{
			"muted":    true,
			"autoplay": false,
			"loop":     false,
			"controls": true,
		},
	)

	viper.SetDefault(
		"swiper.css",
		[]string{
			"static/swiper/swiper-bundle.min.css",
			"static/swiper/swiper.css",
			"static/swiper/gallery.css",
		},
	)

	viper.SetDefault("swiper.mime", []string{"image", "video"})

	viper.SetDefault(
		"swiper.scripts",
		[]string{
			"static/swiper/swiper-bundle.min.js",
		},
	)

	viper.SetDefault(
		"color",
		map[string]string{
			"bg":     "#262626",
			"fg":     "#ffbf00",
			"white":  "#ffffff",
			"black":  "#262626",
			"grey":   "#626262",
			"yellow": "#ffff87",
			"red":    "#ff5f5f",
			"pink":   "#ffafff",
			"cyan":   "#afffff",
			"blue":   "#5fafff",
			"green":  "#afffaf",
			"purple": "#af87ff",
		},
	)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		//home, err := os.UserHomeDir()
		//cobra.CheckErr(err)

		// Search config in home directory with name ".static" (without extension).
		cfgDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Join(cfgDir, "static")
		viper.AddConfigPath(path)
		viper.SetConfigType("toml")
		viper.SetConfigName("config")

		static.UserCfg = os.DirFS(path)
	}

	viper.AutomaticEnv() // read in environment variables that match

	//files, _ := fs.Glob(static.Public, "static/*")
	//fmt.Printf("files %+V\n", files)
	//usr, _ := fs.Glob(static.UserCfg, "swiper/*")
	//fmt.Printf("files %+V\n", usr)
	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		//fmt.Printf("%+V\n", viper.AllSettings())
		//fmt.Printf("%+V\n", viper.Get("gifv.html"))
		//cfile := viper.ConfigFileUsed()
		//err := viper.Unmarshal(&cfg)
		//if err != nil {
		//log.Fatal(err)
		//}
		//cfg, err = static.ParseConfig(cfile)
		//if err != nil {
		//  log.Fatal(err)
		//}
	}
	static.InitTemplates()
}
