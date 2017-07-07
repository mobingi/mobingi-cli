package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mocli",
	Short: "Mobingi API command line interface.",
	Long:  `Command line interface for Mobingi API and services.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("api-version", "v", "v2", "API version")
}
