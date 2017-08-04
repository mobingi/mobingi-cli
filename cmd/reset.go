package cmd

import (
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/debug"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func ResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "reset config to defaults",
		Long:  `Reset all configuration values to default.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := cli.SetDefaultCliConfig()
			err = errors.Wrap(err, "write default config failed")
			debug.ErrorExit(err, 1)
		},
	}
}
