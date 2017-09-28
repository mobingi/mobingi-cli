package cmd

import (
	"fmt"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/credentials"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
	"github.com/spf13/cobra"
)

func userPass(cmd *cobra.Command) *credentials.UserPass {
	userpass := &credentials.UserPass{
		Username: cli.GetCliStringFlag(cmd, "username"),
		Password: cli.GetCliStringFlag(cmd, "password"),
	}

	in, err := userpass.EnsureInput(false)
	if err != nil {
		cli.ErrorExit(err, 1)
	}

	if in[1] {
		fmt.Println("\n") // new line after the password input
	}

	return userpass
}

func ensureUserPass(cmd *cobra.Command, sess *session.Session) *credentials.UserPass {
	var userpass *credentials.UserPass
	if sess.Config.Username == "" {
		userpass = userPass(cmd)
		sess.Config.Username = userpass.Username
		sess.Config.Password = userpass.Password
	}

	return userpass
}
