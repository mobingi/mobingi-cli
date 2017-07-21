package cmd

import (
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/spf13/cobra"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "manage your docker registry",
	Long:  `Manage your docker registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		d.Info("Check `registry --help` for more information on supported subcommands.")
	},
}

func init() {
	rootCmd.AddCommand(registryCmd)
}
