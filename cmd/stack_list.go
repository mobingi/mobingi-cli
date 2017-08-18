package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mobingi/mobingi-cli/client"
	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/pretty"
	"github.com/mobingi/mobingi-cli/pkg/stack"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/private/debug"
	"github.com/spf13/cobra"
)

func StackListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all stacks",
		Long: `List all stacks. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

Valid format values: min (default), text, json, raw

For now, the 'min' format option cannot yet write to a file
using the '--out=[filename]' option. You need to specify either
'text', 'json', or 'raw'.

Examples:

  $ ` + cli.BinName() + ` stack list
  $ ` + cli.BinName() + ` stack list --fmt=text
  $ ` + cli.BinName() + ` stack list --fmt=json --verbose
  $ ` + cli.BinName() + ` stack list --fmt=raw --out=/home/foo/tmp.txt`,
		Run: stackList,
	}

	return cmd
}

func stackList(cmd *cobra.Command, args []string) {
	c := client.NewClient(client.NewApiConfig(cmd))
	body, err := c.AuthGet("/alm/stack")
	d.ErrorExit(err, 1)

	var stacks []stack.ListStack
	err = json.Unmarshal(body, &stacks)
	d.ErrorExit(err, 1)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))

		// write to file option
		f := cli.GetCliStringFlag(cmd, "out")
		if f != "" {
			err = ioutil.WriteFile(f, body, 0644)
			d.ErrorExit(err, 1)
			d.Info(fmt.Sprintf("Output written to %s.", f))
		}
	case "text":
		indent := cli.GetCliIntFlag(cmd, "indent")
		stack.PrintR(os.Stdout, &stacks[0], 0, indent)

		// write to file option
		f := cli.GetCliStringFlag(cmd, "out")
		if f != "" {
			fp, err := os.Create(f)
			d.ErrorExit(err, 1)

			defer fp.Close()
			w := bufio.NewWriter(fp)
			defer w.Flush()
			stack.PrintR(w, &stacks[0], 0, indent)
			d.Info(fmt.Sprintf("Output written to %s.", f))
		}
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		mi, err := json.MarshalIndent(stacks, "", pretty.Indent(indent))
		d.ErrorExit(err, 1)

		fmt.Println(string(mi))

		// write to file option
		f := cli.GetCliStringFlag(cmd, "out")
		if f != "" {
			err = ioutil.WriteFile(f, mi, 0644)
			d.ErrorExit(err, 1)
			d.Info(fmt.Sprintf("Output written to %s.", f))
		}
	default:
		if pfmt == "min" || pfmt == "" {
			w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
			fmt.Fprintf(w, "STACK ID\tSTACK NAME\tPLATFORM\tSTATUS\tREGION\tLAUNCHED\n")
			for _, s := range stacks {
				timestr := s.CreateTime
				t, err := time.Parse(time.RFC3339, s.CreateTime)
				if err == nil {
					timestr = t.Format(time.RFC1123)
				}

				platform := "?"
				if s.Configuration.AWS != "" {
					platform = "AWS"
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					s.StackId,
					s.Nickname,
					platform,
					s.StackStatus,
					s.Configuration.Region,
					timestr)
			}

			w.Flush()
		}
	}
}
