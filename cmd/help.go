package cmd

import (
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/spf13/cobra"
)

func HelpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "help about any command",
		Long: `Help provides help for any command in the application.
Simply type '` + cli.BinName() + ` help [path to command]' for full details.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
}
