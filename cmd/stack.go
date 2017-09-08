package cmd

import "github.com/spf13/cobra"

func StackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stack",
		Short: "manage your stack",
		Long:  `Manage your infrastructure/application stack.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		StackListCmd(),
		StackDescribeCmd(),
		StackCreateCmd(),
		StackUpdateCmd(),
		StackDeleteCmd(),
		StackGetPemCmd(),
	)

	return cmd
}
