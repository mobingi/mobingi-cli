package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mobingi/mobingi-cli/client"
	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/iohelper"
	"github.com/mobingi/mobingi-cli/pkg/stack"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

func StackDescribeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "display stack details",
		Long: `Display stack details. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

Valid format values: min (default), json, raw

Examples:

  $ ` + cmdline.Args0() + ` stack describe --id=58c2297d25645-Y6NSE4VjP-tk
  $ ` + cmdline.Args0() + ` stack describe --id=58c2297d25645-Y6NSE4VjP-tk --fmt=json`,
		Run: describe,
	}

	cmd.Flags().StringP("id", "i", "", "stack id")
	return cmd
}

func describe(cmd *cobra.Command, args []string) {
	var err error
	id := cli.GetCliStringFlag(cmd, "id")
	if id == "" {
		d.ErrorExit("stack id cannot be empty", 1)
	}

	c := client.NewClient(client.NewApiConfig(cmd))
	body, err := c.AuthGet("/alm/stack/" + fmt.Sprintf("%s", id))
	d.ErrorExit(err, 1)

	// we process `--fmt=raw` option first
	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	if pfmt == "raw" {
		fmt.Println(string(body))
		if out != "" {
			err = iohelper.WriteToFile(out, body)
			d.ErrorExit(err, 1)
		}

		return
	}

	// workaround: see description in struct definition
	var stacks1 []stack.DescribeStack1
	var stacks2 []stack.DescribeStack2
	valid := 0
	err = json.Unmarshal(body, &stacks1)
	if err != nil {
		err = json.Unmarshal(body, &stacks2)
		d.ErrorExit(err, 1)
		valid = 2
	} else {
		valid = 1
	}

	switch pfmt {
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)

		// write to file option
		if out != "" {
			err = iohelper.WriteToFile(out, []byte(js))
			d.ErrorExit(err, 1)
		}
	default:
		if pfmt == "min" || pfmt == "" {
			w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
			fmt.Fprintf(w, "INSTANCE ID\tINSTANCE TYPE\tINSTANCE MODEL\tPUBLIC IP\tPRIVATE IP\tSTATUS\n")
			if valid == 1 {
				for _, inst := range stacks1[0].Instances {
					instype := "on-demand"
					if inst.InstanceLifecycle == "spot" {
						instype = inst.InstanceLifecycle
					}

					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
						inst.InstanceId,
						instype,
						inst.InstanceType,
						inst.PublicIpAddress,
						inst.PrivateIpAddress,
						inst.State.Name)
				}
			}

			if valid == 2 {
				for _, inst := range stacks2[0].Instances {
					instype := "on-demand"
					if inst.InstanceLifecycle == "spot" {
						instype = inst.InstanceLifecycle
					}

					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
						inst.InstanceId,
						instype,
						inst.InstanceType,
						inst.PublicIpAddress,
						inst.PrivateIpAddress,
						inst.State.Name)
				}
			}

			w.Flush()
		}
	}
}
