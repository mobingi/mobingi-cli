package cmd

import "github.com/spf13/cobra"

func TemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "manage your ALM templates",
		Long:  `Manage your ALM templates.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		TemplateVersionsListCmd(),
		TemplateDescribeCmd(),
		TemplateCompareCmd(),
	)

	return cmd
}
