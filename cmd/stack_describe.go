package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/iohelper"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
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
	sess, err := clisession()
	d.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.StackDescribeInput{
		StackId: cli.GetCliStringFlag(cmd, "id"),
	}

	resp, body, err := svc.Describe(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

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

	var stacks []alm.DescribeStack
	err = json.Unmarshal(body, &stacks)
	d.ErrorExit(err, 1)

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
			for _, inst := range stacks[0].Instances {
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

			w.Flush()
		}
	}
}
