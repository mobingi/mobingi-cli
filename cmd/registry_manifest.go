package cmd

import (
	"fmt"
	"strings"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/credentials"
	"github.com/mobingilabs/mocli/pkg/registry"
	"github.com/spf13/cobra"
)

func RegistryManifest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "print a tag manifest",
		Long: `Print a tag manifest. At the very least, you only have to provide 'username', 'password',
and image name. Other values will be built based on inputs and command type. Output format is JSON.`,
		Run: manifest,
	}

	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	cmd.Flags().String("scope", "", "scope for authentication")
	cmd.Flags().String("image", "", "image name (format: `image:tag`)")
	return cmd
}

func manifest(cmd *cobra.Command, args []string) {
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

	imagepair := strings.Split(image, ":")
	if len(imagepair) != 2 {
		check.ErrorExit("--image format is `image:tag`", 1)
	}

	if scope == "" {
		scope = fmt.Sprintf("repository:%s/%s:pull", up.Username, imagepair[0])
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

	c := client.NewClient(&client.Config{
		RootUrl:     "https://dockereg2.labs.mobingi.com",
		ApiVersion:  "v2",
		AccessToken: token,
	})

	path := fmt.Sprintf("/%s/%s/manifests/%s", up.Username, imagepair[0], imagepair[1])
	_, body, errs := c.Get(path)
	check.ErrorExit(errs, 1)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	default:
		fmt.Println(string(body))
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
