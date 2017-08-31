package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print the version",
		Long:  `Print the version.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: provide a formal versioning system
			fmt.Println("v0.2.1-beta")
		},
	}
}
