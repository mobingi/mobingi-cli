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
	rootCmd.PersistentFlags().StringP("fmt", "f", "", "output format (values depends on command)")
	rootCmd.PersistentFlags().StringP("out", "o", "", "full file path to write the output")
	rootCmd.PersistentFlags().IntP("indent", "n", 4, "indent padding when fmt is 'text' or 'json'")
}
