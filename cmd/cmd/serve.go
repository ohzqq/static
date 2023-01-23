package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var port string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		fs := http.FileServer(http.Dir(dir))
		http.Handle("/", fs)

		log.Printf("Serving %s on HTTP port: %s\n", dir, port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&port, "port", "p", "3000", "port to serve files on")
}
