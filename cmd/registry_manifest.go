package cmd

import (
	"fmt"
	"strings"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/cli"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/iohelper"
	"github.com/mobingilabs/mocli/pkg/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RegistryManifest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "print a tag manifest",
		Long: `Print a tag manifest. At the very least, you only have to provide 'username', 'password',
and image name. Other values will be built based on inputs and command type. Output format is JSON.

Example:

  $ ` + cli.BinName() + ` registry manifest --username=foo --password=bar --image=hello:latest`,
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
	userpass := userPass(cmd)
	base := viper.GetString("api_url")
	apiver := cli.GetCliStringFlag(cmd, "apiver")
	svc := cli.GetCliStringFlag(cmd, "service")
	scope := cli.GetCliStringFlag(cmd, "scope")
	image := cli.GetCliStringFlag(cmd, "image")
	if image == "" {
		d.ErrorExit("image name cannot be empty", 1)
	}

	pair := strings.Split(image, ":")
	if len(pair) != 2 {
		d.ErrorExit("--image format is `image:tag`", 1)
	}

	if scope == "" {
		scope = fmt.Sprintf("repository:%s/%s:pull", userpass.Username, pair[0])
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

	d.ErrorExit(err, 1)

	c := client.NewClient(&client.Config{
		RootUrl:     viper.GetString("registry_url"),
		ApiVersion:  cli.DockerApiVersion,
		AccessToken: token,
	})

	path := fmt.Sprintf("/%s/%s/manifests/%s", userpass.Username, pair[0], pair[1])
	body, err = c.AuthGet(path)
	d.ErrorExit(err, 1)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	default:
		fmt.Println(string(body))
	}

	out := cli.GetCliStringFlag(cmd, "out")
	if out != "" {
		err = iohelper.WriteToFile(out, body)
		d.ErrorExit(err, 1)
	}
}
