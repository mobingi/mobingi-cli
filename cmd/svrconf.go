package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var svrconfCmd = &cobra.Command{
	Use:   "svrconf",
	Short: "manage your server config file",
	Long:  `Manage your server config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Check `svrconf --help` for more information on supported subcommands.")
	},
}

func init() {
	rootCmd.AddCommand(svrconfCmd)
}
