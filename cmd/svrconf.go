package cmd

import (
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/spf13/cobra"
)

var svrconfCmd = &cobra.Command{
	Use:   "svrconf",
	Short: "manage your server config file",
	Long:  `Manage your server config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		d.Info("Check `svrconf --help` for more information on supported subcommands.")
	},
}

func init() {
	rootCmd.AddCommand(svrconfCmd)
}
