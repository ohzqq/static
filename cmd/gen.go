package cmd

import (
	"fmt"
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
		var err error

		arg, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer arg.Close()
		DirCheck(arg)

		entries, err := WalkDir(args[0])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v\n", entries)
	},
}

func WalkDir(root string) ([][]string, error) {
	var files [][]string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		var dirFiles []string

		if d.IsDir() {
			return nil
		}

		if d.Name() == "meta.toml" {
			dirFiles = append(dirFiles, path)
		}

		if d.Name() == "body.html" {
			dirFiles = append(dirFiles, path)
		}

		files = append(files, dirFiles)

		return nil
	})
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
