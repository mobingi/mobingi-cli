package cmd

import "github.com/spf13/cobra"

func CredentialsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "creds",
		Short: "manage your credentials",
		Long:  `Manage your credentials.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		CredentialsListCmd(),
	)

	return cmd
}
