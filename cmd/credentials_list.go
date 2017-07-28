package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	term "github.com/buger/goterm"
	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/credentials"
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
	check.ErrorExit(err, 1)

	var creds []credentials.VendorCredentials
	err = json.Unmarshal(body, &creds)
	check.ErrorExit(err, 1)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		mi, err := json.MarshalIndent(creds, "", pretty.Indent(indent))
		check.ErrorExit(err, 1)

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
}