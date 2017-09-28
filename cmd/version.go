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
			fmt.Println("v1.0.2")
		},
	}
}
