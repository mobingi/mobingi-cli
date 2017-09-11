package cmd

import "github.com/spf13/cobra"

func RbacCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rbac",
		Short: "manage role based access control features",
		Long:  `Manage your role based access control features.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		RbacDescribeCmd(),
		RbacCreateCmd(),
		RbacDeleteCmd(),
		RbacSampleCmd(),
	)

	return cmd
}
