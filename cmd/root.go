package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mocli",
	Short: "",
	Long:  `Mobingi API command line interface.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("token", "", "access token for API access")
}
