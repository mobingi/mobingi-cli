package cmd

import (
	"fmt"

	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/credentials"
	"github.com/spf13/cobra"
)

func userPass(cmd *cobra.Command) *credentials.UserPass {
	up := &credentials.UserPass{
		Username: cli.GetCliStringFlag(cmd, "username"),
		Password: cli.GetCliStringFlag(cmd, "password"),
	}

	in, err := up.EnsureInput(false)
	if err != nil {
		check.ErrorExit(err, 1)
	}

	if in[1] {
		fmt.Println("\n") // new line after the password input
	}

	return up
}
