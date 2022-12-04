package cmd

import (
	"fmt"
	"idxgen/page"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//var err error

		for _, cat := range cfg.Categories {
			col := cfg.Collection[cat]
			path := filepath.Join(args[0], cat)
			idx := page.NewCollection(path, col.Ext...)
			idx.Collection = col
			idx.Type = cat

			fmt.Printf("%v idx %+V\n", idx.Title(), idx.Files)
			for _, c := range idx.Children {
				fmt.Printf("%v idx.Children %+V\n", c.Title(), c.Files)
			}
		}
	},
}

func GetDirEntries(name string) []os.DirEntry {
	//abs, err := filepath.Abs(name)
	//if err != nil {
	//log.Fatal(err)
	//}
	abs := name
	println(abs)
	entries, err := os.ReadDir(abs)
	if err != nil {
		log.Fatal(err)
	}
	return entries
}

func WalkDir(root string) ([][]string, error) {
	var files [][]string
	var dirFiles []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if d.Name() == "meta.toml" {
			dirFiles = append(dirFiles, path)
		}

		if d.Name() == "body.html" {
			dirFiles = append(dirFiles, path)
		}

		return nil
	})
	files = append(files, dirFiles)
	return files, err
}

func DirCheck(f *os.File) {
	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if !stat.IsDir() {
		log.Fatal("not a dir")
	}
}

func init() {
	rootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
