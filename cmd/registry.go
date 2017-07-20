package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "manage your docker registry",
	Long:  `Manage your docker registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Check `registry --help` for more information on supported subcommands.")
	},
}

func init() {
	rootCmd.AddCommand(registryCmd)
}
