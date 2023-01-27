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
	indexOnly  bool
	builder    = &static.Site{}
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

	// defaults
	viper.SetDefault("build.index_only", true)
	viper.SetDefault("build.profile", "global")
	viper.SetDefault("build.format", "yml")
	viper.SetDefault("build.assets", false)

	viper.SetDefault(
		"global.css",
		[]string{
			"public/css/normalize.css",
			"public/css/milligram.css",
		},
	)

	viper.SetDefault("global.template", "filterableList")

	viper.SetDefault(
		"global.scripts",
		[]string{
			"public/js/alpine.js",
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
		"global.html.audio",
		map[string]any{
			"controls": true,
			"muted":    true,
			"loop":     false,
		},
	)

	viper.SetDefault(
		"swiper.css",
		[]string{
			"public/swiper/swiper-bundle.min.css",
			"public/swiper/swiper-lightbox.css",
			"public/swiper/gallery.css",
		},
	)

	viper.SetDefault("swiper.mime", []string{"image", "video"})

	viper.SetDefault(
		"swiper.scripts",
		[]string{
			"public/swiper/swiper-bundle.min.js",
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

	// flags
	rootCmd.PersistentFlags().BoolP("regen", "g", false, "regenerate index files")
	viper.BindPFlag("build.regen", rootCmd.PersistentFlags().Lookup("regen"))

	rootCmd.PersistentFlags().BoolP("all", "a", false, "list all files in nav")
	viper.BindPFlag("build.all", rootCmd.PersistentFlags().Lookup("all"))

	rootCmd.PersistentFlags().BoolP("recurse", "r", false, "recursive build")
	viper.BindPFlag("build.is_collection", rootCmd.PersistentFlags().Lookup("recurse"))

	rootCmd.PersistentFlags().Bool("no-thumbs", false, "don't generate thumbnails")
	viper.BindPFlag("build.no_thumbs", rootCmd.PersistentFlags().Lookup("no-thumbs"))

	rootCmd.PersistentFlags().BoolVarP(&indexOnly, "index-only", "I", false, "only list index.html files")

	rootCmd.PersistentFlags().StringP("fmt", "f", "", "data export format")
	viper.BindPFlag("build.format", rootCmd.PersistentFlags().Lookup("fmt"))

	rootCmd.PersistentFlags().StringSliceP("ext", "e", []string{}, "glob by ext")
	viper.BindPFlag("build.exts", rootCmd.PersistentFlags().Lookup("ext"))

	rootCmd.PersistentFlags().StringSliceP("mime", "m", []string{}, "glob by mimetype")
	viper.BindPFlag("build.mimes", rootCmd.PersistentFlags().Lookup("mime"))

	rootCmd.PersistentFlags().StringP("profile", "p", "", "build profile")
	viper.BindPFlag("build.profile", rootCmd.PersistentFlags().Lookup("profile"))

	rootCmd.PersistentFlags().StringP("template", "t", "", "set golang template for page")
	viper.BindPFlag("build.template", rootCmd.PersistentFlags().Lookup("template"))

	rootCmd.MarkFlagsMutuallyExclusive("all", "mime")
	rootCmd.MarkFlagsMutuallyExclusive("all", "ext")
	rootCmd.MarkFlagsMutuallyExclusive("all", "index-only")
	rootCmd.MarkFlagsMutuallyExclusive("index-only", "ext")
	rootCmd.MarkFlagsMutuallyExclusive("index-only", "mime")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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
