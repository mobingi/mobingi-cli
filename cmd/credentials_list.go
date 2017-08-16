package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/credentials"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/pretty"
	"github.com/spf13/cobra"
)

func CredentialsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list your credentials",
		Long: `List your credentials.

Supported formats: raw, json

Examples:

  $ ` + cli.BinName() + ` creds list
  $ ` + cli.BinName() + ` creds list --fmt=json
  $ ` + cli.BinName() + ` creds list --fmt=raw`,
		Run: credsList,
	}

	cmd.Flags().String("vendor", "aws", "credentials vendor")
	return cmd
}

func credsList(cmd *cobra.Command, args []string) {
	vendor := cli.GetCliStringFlag(cmd, "vendor")
	c := client.NewClient(client.NewApiConfig(cmd))
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
		w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
		fmt.Fprintf(w, "VENDOR\tID\tACCOUNT\tLAST MODIFIED\n")
		for _, cred := range creds {
			timestr := cred.LastModified
			t, err := time.Parse(time.RFC3339, cred.LastModified)
			if err == nil {
				timestr = t.Format(time.RFC1123)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				vendor,
				cred.Id,
				cred.Account,
				timestr)
		}

		w.Flush()
	}
}

func getCredsList(cmd *cobra.Command) ([]credentials.VendorCredentials, error) {
	return nil, nil
}
