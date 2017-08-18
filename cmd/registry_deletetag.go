package cmd

import (
	"fmt"
	"strings"

	"github.com/mobingi/mobingi-cli/client"
	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/registry"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/private/debug"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RegistryDeleteTag() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete an image tag",
		Long: `Delete an image tag. At the very least, you only have to provide 'username', 'password',
and image name. Other values will be built based on inputs and command type.

Example:

  $ ` + cli.BinName() + ` registry delete --username=foo --password=bar --image=hello:latest`,
		Run: deleteTag,
	}

	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	cmd.Flags().String("scope", "", "scope for authentication")
	cmd.Flags().String("image", "", "image name (format: `image:tag`)")
	return cmd
}

func deleteTag(cmd *cobra.Command, args []string) {
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

	tp := &registry.TokenParams{
		Base:       base,
		ApiVersion: apiver,
		TokenCreds: &registry.TokenCredentials{
			UserPass: userpass,
			Service:  svc,
			Scope:    scope,
		},
	}

	// request token for get manifest (pull)
	_, token, err := registry.GetRegistryToken(tp)
	d.ErrorExit(err, 1)

	rurl := viper.GetString("registry_url")
	c := client.NewClient(&client.Config{
		RootUrl:     rurl,
		ApiVersion:  cli.DockerApiVersion,
		AccessToken: token,
	})

	// get manifest to get tag digest
	path := fmt.Sprintf("/%s/%s/manifests/%s", userpass.Username, pair[0], pair[1])
	digest, err := c.GetTagDigest(path)
	d.ErrorExit(err, 1)

	// new token for delete
	scope = fmt.Sprintf("repository:%s/%s:*", userpass.Username, pair[0])
	tp.TokenCreds.Scope = scope
	_, token, err = registry.GetRegistryToken(tp)
	d.ErrorExit(err, 1)

	c2 := client.NewClient(&client.Config{
		RootUrl:     rurl,
		ApiVersion:  cli.DockerApiVersion,
		AccessToken: token,
	})

	path = fmt.Sprintf("/%s/%s/manifests/%s", userpass.Username, pair[0], digest)
	_, err = c2.AuthDel(path)
	d.ErrorExit(err, 1)

	d.Info(fmt.Sprintf("Tag '%s:%s' deleted.", pair[0], pair[1]))
}
