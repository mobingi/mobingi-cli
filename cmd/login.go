package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "",
	Long:  `Placeholder for the documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("login here")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().String("client-id", "", "client id")
	loginCmd.Flags().String("client-secret", "", "client secret")
	loginCmd.Flags().String("grant-type", "client_credentials", "grant type")
}
