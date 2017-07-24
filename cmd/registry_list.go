package cmd

import "github.com/spf13/cobra"

func RegistryList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list images and tags",
		Long:  `List images and tags.`,
		Run:   rlist,
	}

	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	cmd.Flags().String("scope", "", "scope for authentication")
	return cmd
}

func rlist(cmd *cobra.Command, args []string) {
}
