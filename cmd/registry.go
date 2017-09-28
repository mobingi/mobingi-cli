package cmd

import "github.com/spf13/cobra"

func RegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "manage your Mobingi Docker Registry",
		Long:  `Manage your Mobingi Docker Registry.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		RegistryToken(),
		RegistryTagsList(),
		RegistryCatalog(),
		DescribeImageCmd(),
		RegistryManifest(),
		RegistryDeleteTag(),
	)

	return cmd
}
