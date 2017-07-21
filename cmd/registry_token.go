package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/constants"
	"github.com/mobingilabs/mocli/pkg/credentials"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/registry"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "get registry token",
	Long: `Get registry token. This command supports '--fmt=raw' option. By default,
it will only print the token value.

If you want to save the token for other registry-related commands,
use the '--out=home' option.`,
	Run: token,
}

func init() {
	registryCmd.AddCommand(tokenCmd)
	tokenCmd.Flags().String("account", "", "subuser name")
	tokenCmd.Flags().String("username", "", "username (account subuser)")
	tokenCmd.Flags().String("password", "", "password (account subuser)")
	tokenCmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	tokenCmd.Flags().String("scope", "", "scope for authentication")
}

func token(cmd *cobra.Command, args []string) {
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

	base := cli.GetCliStringFlag(cmd, "url")
	apiver := cli.GetCliStringFlag(cmd, "apiver")
	acct := cli.GetCliStringFlag(cmd, "account")
	svc := cli.GetCliStringFlag(cmd, "service")
	scope := cli.GetCliStringFlag(cmd, "scope")
	if acct == "" {
		acct = up.Username
	}

	body, token, err := registry.GetRegistryToken(&registry.TokenParams{
		Base:       base,
		ApiVersion: apiver,
		TokenCreds: &registry.TokenCredentials{
			UserPass: up,
			Account:  acct,
			Service:  svc,
			Scope:    scope,
		},
	}, false)

	if err != nil {
		check.ErrorExit(err, 1)
	}

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
	default:
		d.Info("token:", token)
	}

	out := cli.GetCliStringFlag(cmd, "out")
	if out != "" {
		switch out {
		case "home":
			err = credentials.SaveRegistryToken(token)
			if err != nil {
				check.ErrorExit(err, 1)
			}

			hd := credentials.CredFolder(false)
			rf := filepath.Join(hd, constants.REGTOKEN_FILE)
			d.Info(fmt.Sprintf("output written to %s", rf))
		default:
			d.Error("should set '--out=home' option")
		}
	}
}
