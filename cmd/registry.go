package cmd

import "github.com/spf13/cobra"

func RegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "manage your Mobingi docker registry",
		Long:  `Manage your Mobingi docker registry.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		RegistryToken(),
		RegistryTagsList(),
		RegistryCatalog(),
		RegistryManifest(),
		RegistryDeleteTag(),
	)

	return cmd
}
