package cmd

import (
	"github.com/spf13/cobra"
)

func RegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "manage your docker registry",
		Long:  `Manage your docker registry.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		RegistryToken(),
		RegistryList(),
	)

	return cmd
}
