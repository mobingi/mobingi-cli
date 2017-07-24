package cmd

import (
	"log"
	"os"

	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mocli",
	Short: "Mobingi API command line interface.",
	Long:  `Command line interface for Mobingi API and services.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		d.Info("prerun")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		d.Info("postrun")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("token", "", "access token")
	rootCmd.PersistentFlags().String("url", "https://apidev.mobingi.com", "API endpoint url")
	rootCmd.PersistentFlags().String("apiver", "v2", "API version")
	rootCmd.PersistentFlags().StringP("fmt", "f", "", "output format (values depends on command)")
	rootCmd.PersistentFlags().StringP("out", "o", "", "full file path to write the output")
	rootCmd.PersistentFlags().IntP("indent", "n", 4, "indent padding when fmt is 'text' or 'json'")
}
