package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	term "github.com/buger/goterm"
	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/iohelper"
	"github.com/mobingilabs/mocli/pkg/pretty"
	"github.com/mobingilabs/mocli/pkg/stack"
	"github.com/spf13/cobra"
)

func StackDescribeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "display stack details",
		Long: `Display stack details. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

Valid format values: text (default), json, raw, min`,
		Run: describe,
	}

	cmd.Flags().StringP("id", "i", "", "stack id")
	return cmd
}

func describe(cmd *cobra.Command, args []string) {
	var err error
	id := cli.GetCliStringFlag(cmd, "id")
	if id == "" {
		check.ErrorExit("stack id cannot be empty", 1)
	}

	c := client.NewClient(client.NewApiConfig(cmd))
	resp, body, errs := c.Get("/alm/stack/" + fmt.Sprintf("%s", id))
	check.ErrorExit(errs, 1)

	// we process `--fmt=raw` option first
	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	if pfmt == "raw" {
		fmt.Println(string(body))
		if out != "" {
			err = iohelper.WriteToFile(out, body)
			check.ErrorExit(err, 1)
		}

		return
	}

	// workaround: see description in struct definition
	var ptr interface{}  // pointer to 1st element of slice
	var sptr interface{} // pointer to the whole slice
	var stacks1 []stack.DescribeStack1
	var stacks2 []stack.DescribeStack2
	valid := 0
	err = json.Unmarshal(body, &stacks1)
	if err != nil {
		err = json.Unmarshal(body, &stacks2)
		if err != nil {
			serr := check.ResponseError(resp, body)
			check.ErrorExit(serr, 1)
			check.ErrorExit(err, 1)
		} else {
			ptr = &stacks2[0]
			sptr = stacks2
			valid = 2
		}
	} else {
		ptr = &stacks1[0]
		sptr = stacks1
		valid = 1
	}

	switch pfmt {
	case "min":
		stbl := term.NewTable(0, 10, 5, ' ', 0)
		fmt.Fprintf(stbl, "INSTANCE ID\tINSTANCE TYPE\tPUBLIC IP\tPRIVATE IP\tSTATUS\n")
		if valid == 1 {
			for _, inst := range stacks1[0].Instances {
				fmt.Fprintf(stbl, "%s\t%s\t%s\t%s\t%s\n",
					inst.InstanceId,
					inst.InstanceType,
					inst.PublicIpAddress,
					inst.PrivateIpAddress,
					inst.State.Name)
			}
		}

		if valid == 2 {
			for _, inst := range stacks2[0].Instances {
				fmt.Fprintf(stbl, "%s\t%s\t%s\t%s\t%s\n",
					inst.InstanceId,
					inst.InstanceType,
					inst.PublicIpAddress,
					inst.PrivateIpAddress,
					inst.State.Name)
			}
		}

		term.Print(stbl)
		term.Flush()
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		mi, err := json.MarshalIndent(sptr, "", pretty.Indent(indent))
		check.ErrorExit(err, 1)

		fmt.Println(string(mi))

		// write to file option
		if out != "" {
			err = iohelper.WriteToFile(out, mi)
			check.ErrorExit(err, 1)
		}
	default:
		if pfmt == "text" || pfmt == "" {
			indent := cli.GetCliIntFlag(cmd, "indent")
			stack.PrintR(os.Stdout, ptr, 0, indent)

			// write to file option
			if out != "" {
				fp, err := os.Create(out)
				check.ErrorExit(err, 1)

				defer fp.Close()
				w := bufio.NewWriter(fp)
				defer w.Flush()
				stack.PrintR(w, ptr, 0, indent)
				d.Info(fmt.Sprintf("output written to %s", out))
			}
		}
	}
}
