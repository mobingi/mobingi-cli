package cmd

import (
	"encoding/json"
	"fmt"

	term "github.com/buger/goterm"
	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/credentials"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/registry"
	"github.com/spf13/cobra"
)

func RegistryCatalog() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "list catalog images",
		Long:  `List catalog images.`,
		Run:   catalog,
	}

	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("account", "", "subuser name")
	cmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	cmd.Flags().String("scope", "", "scope for authentication")
	cmd.Flags().String("image", "", "image name to query")
	return cmd
}

func catalog(cmd *cobra.Command, args []string) {
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
	svc := cli.GetCliStringFlag(cmd, "service")
	scope := cli.GetCliStringFlag(cmd, "scope")
	image := cli.GetCliStringFlag(cmd, "image")
	if image == "" {
		check.ErrorExit("image name cannot be empty", 1)
	}

	if scope == "" {
		scope = fmt.Sprintf("repository:%s/%s:pull", up.Username, image)
	}

	body, token, err := registry.GetRegistryToken(&registry.TokenParams{
		Base:       base,
		ApiVersion: apiver,
		TokenCreds: &registry.TokenCredentials{
			UserPass: up,
			Service:  svc,
			Scope:    scope,
		},
	}, false)

	if err != nil {
		check.ErrorExit(err, 1)
	}

	c := client.NewClient(&client.Config{
		RootUrl:     "https://dockereg2.labs.mobingi.com",
		ApiVersion:  "v2",
		AccessToken: token,
	})

	path := fmt.Sprintf("/%s/%s/tags/list", up.Username, image)
	_, body, errs := c.Get(path)
	check.ErrorExit(errs, 1)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
	default:
		if d.Verbose {
			d.Info("Token used:", token)
		}

		var t tags
		err = json.Unmarshal(body, &t)
		if err != nil {
			check.ErrorExit(err, 1)
		}

		stbl := term.NewTable(0, 10, 5, ' ', 0)
		fmt.Fprintf(stbl, "IMAGE\tTAG\n")
		for _, v := range t.Tags {
			fmt.Fprintf(stbl, "%s\t%s\n", t.Name, v)
		}

		term.Print(stbl)
		term.Flush()
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
