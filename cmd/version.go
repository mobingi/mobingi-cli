package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var verCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version",
	Long:  `Print the version.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: provide a formal versioning system
		fmt.Println("v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(verCmd)
}
