package cmd

import (
	"fmt"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/registry"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/spf13/cobra"
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
	service := cli.GetCliStringFlag(cmd, "service")
	scope := cli.GetCliStringFlag(cmd, "scope")
	sess, err := clisession()
	cli.ErrorExit(err, 1)

	ensureUserPass(cmd, sess)
	svc := registry.New(sess)
	in := &registry.GetRegistryTokenInput{
		Service: service,
		Scope:   scope,
	}

	resp, body, token, err := svc.GetRegistryToken(in)
	cli.ErrorExit(err, 1)
	exitOn401(resp)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
	default:
		d.Info("token:", token)
	}
}
