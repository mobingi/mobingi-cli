package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	term "github.com/buger/goterm"
	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/pretty"
	"github.com/mobingilabs/mocli/pkg/stack"
	"github.com/spf13/cobra"
)

func StackListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all stacks",
		Long: `List all stacks. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

Valid format values: min (default), text, json

For now, the 'min' format option cannot yet write to a file
using the '--out=[filename]' option. You need to specify either
'text' or 'json'.`,
		Run: slist,
	}

	return cmd
}

func slist(cmd *cobra.Command, args []string) {
	c := client.NewClient(client.NewConfig(cmd))
	resp, body, errs := c.Get("/alm/stack")
	check.ErrorExit(errs, 1)

	var stacks []stack.ListStack
	err := json.Unmarshal(body, &stacks)
	if err != nil {
		serr := check.ResponseError(resp, body)
		check.ErrorExit(serr, 1)
	}

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "text":
		indent := cli.GetCliIntFlag(cmd, "indent")
		stack.PrintR(os.Stdout, &stacks[0], 0, indent)

		// write to file option
		f := cli.GetCliStringFlag(cmd, "out")
		if f != "" {
			fp, err := os.Create(f)
			check.ErrorExit(err, 1)

			defer fp.Close()
			w := bufio.NewWriter(fp)
			defer w.Flush()
			stack.PrintR(w, &stacks[0], 0, indent)
			d.Info(fmt.Sprintf("Output written to %s.", f))
		}
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		mi, err := json.MarshalIndent(stacks, "", pretty.Indent(indent))
		check.ErrorExit(err, 1)

		fmt.Println(string(mi))

		// write to file option
		f := cli.GetCliStringFlag(cmd, "out")
		if f != "" {
			err = ioutil.WriteFile(f, mi, 0644)
			check.ErrorExit(err, 1)
			d.Info(fmt.Sprintf("Output written to %s.", f))
		}
	default:
		if pfmt == "min" || pfmt == "" {
			stbl := term.NewTable(0, 10, 5, ' ', 0)
			fmt.Fprintf(stbl, "STACK ID\tSTACK NAME\tPLATFORM\tSTATUS\tREGION\tLAUNCHED\n")
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

				fmt.Fprintf(stbl, "%s\t%s\t%s\t%s\t%s\t%s\n",
					s.StackId,
					s.Nickname,
					platform,
					s.StackStatus,
					s.Configuration.Region,
					timestr)
			}

			term.Print(stbl)
			term.Flush()
		}
	}
}
