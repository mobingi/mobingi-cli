package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mobingi/mobingi-cli/client"
	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/credentials"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	"github.com/spf13/cobra"
)

func CredentialsAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add credentials",
		Long: `Add credentials.

Examples:

  $ ` + cmdline.Args0() + ` creds list
  $ ` + cmdline.Args0() + ` creds list --fmt=json`,
		Run: credsAdd,
	}

	cmd.Flags().String("vendor", "aws", "credentials vendor")
	cmd.Flags().String("name", "", "credentials name")
	cmd.Flags().String("key", "", "client key id")
	cmd.Flags().String("secret", "", "client secret")
	return cmd
}

func credsAdd(cmd *cobra.Command, args []string) {
	vendor := cli.GetCliStringFlag(cmd, "vendor")
	name := cli.GetCliStringFlag(cmd, "name")
	key := cli.GetCliStringFlag(cmd, "key")
	secret := cli.GetCliStringFlag(cmd, "secret")

	payload, err := json.Marshal(&credentials.AddVendorCredentials{
		Credentials: credentials.AWSCredentials{
			Name:   name,
			KeyId:  key,
			Secret: secret,
		},
	})

	cli.ErrorExit(err, 1)
	fmt.Println(string(payload))

	c := client.NewClient(client.NewApiConfig(cmd))
	body, err := c.AuthPost("/credentials/"+vendor, payload)
	cli.ErrorExit(err, 1)

	fmt.Println(string(body))
}
