package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/credentials"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/spf13/cobra"
)

func CredentialsAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add credentials",
		Long: `Add credentials.

Examples:

  $ ` + cli.BinName() + ` creds list
  $ ` + cli.BinName() + ` creds list --fmt=json`,
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

	d.ErrorExit(err, 1)
	fmt.Println(string(payload))

	c := client.NewClient(client.NewApiConfig(cmd))
	body, err := c.AuthPost("/credentials/"+vendor, payload)
	d.ErrorExit(err, 1)

	fmt.Println(string(body))

	/*
		body, err := c.AuthGet("/credentials/" + vendor)
		d.ErrorExit(err, 1)

		var creds []credentials.VendorCredentials
		err = json.Unmarshal(body, &creds)
		d.ErrorExit(err, 1)

		pfmt := cli.GetCliStringFlag(cmd, "fmt")
		switch pfmt {
		case "raw":
			fmt.Println(string(body))
		case "json":
			indent := cli.GetCliIntFlag(cmd, "indent")
			mi, err := json.MarshalIndent(creds, "", pretty.Indent(indent))
			d.ErrorExit(err, 1)

			fmt.Println(string(mi))
		default:
			stbl := term.NewTable(0, 10, 5, ' ', 0)
			fmt.Fprintf(stbl, "VENDOR\tID\tACCOUNT\tLAST MODIFIED\n")
			for _, cred := range creds {
				timestr := cred.LastModified
				t, err := time.Parse(time.RFC3339, cred.LastModified)
				if err == nil {
					timestr = t.Format(time.RFC1123)
				}

				fmt.Fprintf(stbl, "%s\t%s\t%s\t%s\n",
					vendor,
					cred.Id,
					cred.Account,
					timestr)
			}

			term.Print(stbl)
			term.Flush()
		}
	*/
}
