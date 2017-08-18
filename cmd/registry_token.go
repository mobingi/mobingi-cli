package cmd

import (
	"fmt"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/registry"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/private/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/private/debug"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RegistryToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "get registry token",
		Long: `Get registry token. This command supports '--fmt=raw' option. By default,
it will only print the token value.

Example:

  $ ` + cmdline.Args0() + ` registry token \
      --username=foo \
      --password=bar \
      --scope="repository:foo/sample:pull"`,
		Run: token,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	cmd.Flags().String("scope", "", "scope for authentication")
	return cmd
}

func token(cmd *cobra.Command, args []string) {
	userpass := userPass(cmd)
	base := viper.GetString("api_url")
	apiver := cli.GetCliStringFlag(cmd, "apiver")
	svc := cli.GetCliStringFlag(cmd, "service")
	scope := cli.GetCliStringFlag(cmd, "scope")

	body, token, err := registry.GetRegistryToken(
		&registry.TokenParams{
			Base:       base,
			ApiVersion: apiver,
			TokenCreds: &registry.TokenCredentials{
				UserPass: userpass,
				Service:  svc,
				Scope:    scope,
			},
		},
	)

	d.ErrorExit(err, 1)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
	default:
		d.Info("token:", token)
	}
}
