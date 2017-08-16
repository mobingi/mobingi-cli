package cmd

import (
	"fmt"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/credentials"
	d "github.com/mobingi/mobingi-cli/pkg/debug"
	"github.com/spf13/cobra"
)

func userPass(cmd *cobra.Command) *credentials.UserPass {
	userpass := &credentials.UserPass{
		Username: cli.GetCliStringFlag(cmd, "username"),
		Password: cli.GetCliStringFlag(cmd, "password"),
	}

	in, err := userpass.EnsureInput(false)
	if err != nil {
		d.ErrorExit(err, 1)
	}

	if in[1] {
		fmt.Println("\n") // new line after the password input
	}

	return userpass
}
