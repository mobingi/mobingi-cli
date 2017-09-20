package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/registry"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/spf13/cobra"
)

func RegistryManifest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "print a tag manifest",
		Long: `Print a tag manifest. At the very least, you only have to provide 'username', 'password',
and image name. Other values will be built based on inputs and command type. Output format is JSON.

Example:

  $ ` + cmdline.Args0() + ` registry manifest --image=hello:latest`,
		Run: manifest,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("username", "", "username (account subuser)")
	cmd.Flags().String("password", "", "password (account subuser)")
	cmd.Flags().String("service", "Mobingi Docker Registry", "service for authentication")
	cmd.Flags().String("scope", "", "scope for authentication")
	cmd.Flags().String("image", "", "image name (format: `image:tag`)")
	return cmd
}

func manifest(cmd *cobra.Command, args []string) {
	service := cli.GetCliStringFlag(cmd, "service")
	scope := cli.GetCliStringFlag(cmd, "scope")
	image := cli.GetCliStringFlag(cmd, "image")
	if image == "" {
		d.ErrorExit("image name cannot be empty", 1)
	}

	pair := strings.Split(image, ":")
	if len(pair) != 2 {
		d.ErrorExit("--image format is `image:tag`", 1)
	}

	sess, err := clisession()
	d.ErrorExit(err, 1)

	ensureUserPass(cmd, sess)
	svc := registry.New(sess)
	in := &registry.GetTagManifestInput{
		Service: service,
		Scope:   scope,
		Image:   pair[0],
		Tag:     pair[1],
	}

	resp, body, err := svc.GetTagManifest(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	default:
		fmt.Println(string(body))
	}

	out := cli.GetCliStringFlag(cmd, "out")
	if out != "" {
		err = ioutil.WriteFile(out, body, 0644)
		d.ErrorExit(err, 1)
	}
}
