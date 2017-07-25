package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/constants"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/registry"
	"github.com/spf13/cobra"
)

func RegistryDeleteTag() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete an image tag",
		Long: `Delete an image tag. At the very least, you only have to provide 'username', 'password',
and image name. Other values will be built based on inputs and command type.

Example:

  $ mocli registry delete --username=foo --password=bar --image=hello:latest`,
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
	base := BaseApiUrl(cmd)
	apiver := cli.GetCliStringFlag(cmd, "apiver")
	svc := cli.GetCliStringFlag(cmd, "service")
	scope := cli.GetCliStringFlag(cmd, "scope")
	image := cli.GetCliStringFlag(cmd, "image")
	if image == "" {
		check.ErrorExit("image name cannot be empty", 1)
	}

	pair := strings.Split(image, ":")
	if len(pair) != 2 {
		check.ErrorExit("--image format is `image:tag`", 1)
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
	if err != nil {
		check.ErrorExit(err, 1)
	}

	rurl := constants.PROD_REG_BASE
	if check.IsDevMode() {
		rurl = constants.DEV_REG_BASE
	}

	c := client.NewClient(&client.Config{
		RootUrl:     rurl,
		ApiVersion:  "v2",
		AccessToken: token,
	})

	// get manifest to get tag digest
	path := fmt.Sprintf("/%s/%s/manifests/%s", userpass.Username, pair[0], pair[1])
	xhdrs := http.Header{
		"Accept": {"application/vnd.docker.distribution.manifest.v2+json"},
	}

	hdrs, err := c.GetHeaders(path, url.Values{}, xhdrs)
	check.ErrorExit(err, 1)

	var digest string
	for n, h := range hdrs {
		if n == "Etag" {
			digest = h[0]
			digest = strings.TrimSuffix(strings.TrimPrefix(digest, "\""), "\"")
		}
	}

	if digest == "" {
		check.ErrorExit("digest not found", 1)
	}

	scope = fmt.Sprintf("repository:%s/%s:*", userpass.Username, pair[0])
	tp.TokenCreds.Scope = scope
	_, token, err = registry.GetRegistryToken(tp)
	if err != nil {
		check.ErrorExit(err, 1)
	}

	c2 := client.NewClient(&client.Config{
		RootUrl:     rurl,
		ApiVersion:  "v2",
		AccessToken: token,
	})

	path = fmt.Sprintf("/%s/%s/manifests/%s", userpass.Username, pair[0], digest)
	_, err = c2.Del(path, url.Values{})
	check.ErrorExit(err, 1)
	d.Info(fmt.Sprintf("Tag '%s:%s' deleted.", pair[0], pair[1]))
}
