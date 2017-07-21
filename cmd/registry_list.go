package cmd

import (
	"github.com/spf13/cobra"
)

var rlistCmd = &cobra.Command{
	Use:   "list",
	Short: "list images and tags",
	Long:  `List images and tags.`,
	Run:   rlist,
}

func init() {
	registryCmd.AddCommand(rlistCmd)
	rlistCmd.Flags().String("username", "", "username (account subuser)")
	rlistCmd.Flags().String("password", "", "password (account subuser)")
	rlistCmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	rlistCmd.Flags().String("scope", "", "scope for authentication")
}

func rlist(cmd *cobra.Command, args []string) {
}
