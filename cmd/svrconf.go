package cmd

import "github.com/spf13/cobra"

func ServerConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "svrconf",
		Short: "manage your server config file",
		Long:  `Manage your server config file.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		ServerConfigShowCmd(),
		ServerConfigUpdateCmd(),
	)

	return cmd
}
