package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "local-build"

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print the version",
		Long:  `Print the version.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}
}
