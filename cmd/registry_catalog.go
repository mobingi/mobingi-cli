package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/cli/confmap"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type catalog struct {
	Repositories []string `json:"repositories"`
}

func RegistryCatalog() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "list catalog images",
		Long: `List catalog images. Note that this command will probably
take some time to complete.

Example:

  $ ` + cli.BinName() + ` registry catalog --username=foo --password=bar`,
		Run: printCatalog,
	}

	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	cmd.Flags().String("scope", "", "scope for authentication")
	return cmd
}

func printCatalog(cmd *cobra.Command, args []string) {
	userpass := userPass(cmd)
	base := viper.GetString("api_url")
	apiver := cli.GetCliStringFlag(cmd, "apiver")
	svc := cli.GetCliStringFlag(cmd, "service")
	scope := cli.GetCliStringFlag(cmd, "scope")
	if scope == "" {
		scope = "registry:catalog:*"
	}

	body, token, err := registry.GetRegistryToken(&registry.TokenParams{
		Base:       base,
		ApiVersion: apiver,
		TokenCreds: &registry.TokenCredentials{
			UserPass: userpass,
			Service:  svc,
			Scope:    scope,
		},
	})

	if err != nil {
		d.ErrorExit(err, 1)
	}

	c := client.NewClient(&client.Config{
		RootUrl:     viper.GetString("registry_url"),
		ApiVersion:  cli.DockerApiVersion,
		AccessToken: token,
	})

	body, err = c.AuthGet("/_catalog")
	d.ErrorExit(err, 1)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
	default:
		if viper.GetBool(confmap.ConfigKey("verbose")) {
			d.Info("[TOKEN USED]", token)
		}

		var ct catalog
		err = json.Unmarshal(body, &ct)
		d.ErrorExit(err, 1)

		for _, v := range ct.Repositories {
			fmt.Println(v)
		}
	}
}
