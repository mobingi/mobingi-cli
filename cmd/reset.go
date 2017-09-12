package cmd

import (
	"os"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/nativestore"
	"github.com/spf13/cobra"
)

func ResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "reset config to defaults",
		Long:  `Reset all configuration values to default.`,
		Run: func(cmd *cobra.Command, args []string) {
			var failed bool
			err := cli.SetDefaultCliConfig()
			if err != nil {
				d.Error(err)
				failed = true
			}

			err = nativestore.Del(cli.CliLabel, cli.CliUrl)
			if err != nil {
				d.Error(err)
				failed = true
			}

			if failed {
				os.Exit(1)
			}
		},
	}
}
