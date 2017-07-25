package cmd

import (
	"fmt"

	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/constants"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/registry"
	"github.com/spf13/cobra"
)

func RegistryToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "get registry token",
		Long: `Get registry token. This command supports '--fmt=raw' option. By default,
it will only print the token value.`,
		Run: token,
	}

	// cmd.Flags().String("account", "", "subuser name")
	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	cmd.Flags().String("scope", "", "scope for authentication")
	return cmd
}

func token(cmd *cobra.Command, args []string) {
	up := userPass(cmd)
	base := cli.GetCliStringFlag(cmd, "url")
	apiver := cli.GetCliStringFlag(cmd, "apiver")
	svc := cli.GetCliStringFlag(cmd, "service")
	scope := cli.GetCliStringFlag(cmd, "scope")
	if base == "" {
		if check.IsDevMode() {
			base = constants.DEV_API_BASE
		} else {
			base = constants.PROD_API_BASE
		}
	}

	body, token, err := registry.GetRegistryToken(&registry.TokenParams{
		Base:       base,
		ApiVersion: apiver,
		TokenCreds: &registry.TokenCredentials{
			UserPass: up,
			Service:  svc,
			Scope:    scope,
		},
	})

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

	/*
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
	*/
}
