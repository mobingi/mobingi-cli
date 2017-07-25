package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/constants"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/registry"
	"github.com/spf13/cobra"
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

  $ mocli registry catalog --username=foo --password=bar`,
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
	base := BaseApiUrl(cmd)
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
		check.ErrorExit(err, 1)
	}

	c := client.NewClient(&client.Config{
		RootUrl:     BaseRegUrl(cmd),
		ApiVersion:  constants.DOCKER_API_VER,
		AccessToken: token,
	})

	body, err = c.Get("/_catalog", url.Values{}, http.Header{})
	check.ErrorExit(err, 1)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
	default:
		if d.Verbose {
			d.Info("Token used:", token)
		}

		var ct catalog
		err = json.Unmarshal(body, &ct)
		check.ErrorExit(err, 1)

		for _, v := range ct.Repositories {
			fmt.Println(v)
		}
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
