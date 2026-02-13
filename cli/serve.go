package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zgiai/kest-api/launcher"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Kest Platform (API + Web Console)",
	Long:  `Run the collaborative platform. This starts the backend API and serves the embedded Web UI.`,
	Run: func(cmd *cobra.Command, args []string) {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		launcher.Start(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
