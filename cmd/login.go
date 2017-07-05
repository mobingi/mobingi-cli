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
	loginCmd.Flags().StringP("client-id", "i", "", "client id")
	loginCmd.Flags().StringP("client-secret", "s", "", "client secret")
	loginCmd.Flags().StringP("grant-type", "g", "client_credentials", "grant type")
}
