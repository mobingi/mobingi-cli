package cmd

import (
	"github.com/spf13/cobra"
)

var svrconfCmd = &cobra.Command{
	Use:   "svrconf",
	Short: "manage your server config file",
	Long:  `Manage your server config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(svrconfCmd)
}
