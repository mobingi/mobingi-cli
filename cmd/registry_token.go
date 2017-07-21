package cmd

import (
	"fmt"

	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/registry"
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "get registry token",
	Long: `Get registry token. This command supports '--fmt=raw' option. By default,
it will only print the token value.`,
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
	passin := false
	user := util.GetCliStringFlag(cmd, "username")
	pass := util.GetCliStringFlag(cmd, "password")
	if user == "" {
		user = util.Username()
	}

	if user == "" {
		util.CheckErrorExit("username cannot be empty", 1)
	}

	if pass == "" {
		pass = util.Password()
		passin = true
	}

	if pass == "" {
		util.CheckErrorExit("password cannot be empty", 1)
	}

	if passin {
		fmt.Println("\n") // new line after the password input
	}

	base := util.GetCliStringFlag(cmd, "url")
	apiver := util.GetCliStringFlag(cmd, "apiver")
	acct := util.GetCliStringFlag(cmd, "account")
	svc := util.GetCliStringFlag(cmd, "service")
	scope := util.GetCliStringFlag(cmd, "scope")
	if acct == "" {
		acct = user
	}

	body, token, err := registry.GetRegistryToken(&registry.TokenParams{
		Base:       base,
		ApiVersion: apiver,
		Username:   user,
		Password:   pass,
		Account:    acct,
		Service:    svc,
		Scope:      scope,
	})

	if err != nil {
		util.CheckErrorExit(err, 1)
	}

	pfmt := util.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
	default:
		d.Info("token:", token)

	}
}
