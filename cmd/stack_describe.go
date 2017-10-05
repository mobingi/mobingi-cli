package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mobingi/mobingi-cli/pkg/cli"
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

	cmd.Flags().String("id", "", "stack id")
	return cmd
}

func describe(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	cli.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.StackDescribeInput{
		StackId: cli.GetCliStringFlag(cmd, "id"),
	}

	resp, body, err := svc.Describe(in)
	cli.ErrorExit(err, 1)
	exitOn401(resp)

	// we process `--fmt=raw` option first
	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
		if out != "" {
			err = ioutil.WriteFile(out, body, 0644)
			cli.ErrorExit(err, 1)
		}
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)

		// write to file option
		if out != "" {
			err = ioutil.WriteFile(out, []byte(js), 0644)
			cli.ErrorExit(err, 1)
		}
	default:
		if pfmt == "min" || pfmt == "" {
			var stacks []alm.DescribeStack
			var stack alm.DescribeStack

			switch sess.Config.ApiVersion {
			case 3:
				err = json.Unmarshal(body, &stack)
				cli.ErrorExit(err, 1)
			default:
				err = json.Unmarshal(body, &stacks)
				cli.ErrorExit(err, 1)
				stack = stacks[0]
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
			fmt.Fprintf(w, "INSTANCE ID\tINSTANCE TYPE\tINSTANCE MODEL\tPUBLIC DNS\tPUBLIC IP\tPRIVATE IP\tSTATUS\n")
			for _, inst := range stack.Instances {
				instype := "on-demand"
				if inst.InstanceLifecycle == "spot" {
					instype = inst.InstanceLifecycle
				}

				pubip := string(inst.PublicIpAddress)
				pubip = strings.TrimLeft(pubip, "\"")
				pubip = strings.TrimRight(pubip, "\"")

				// try if ip is json (alicloud)
				type pubip_t struct {
					IpAddress []string
				}

				var ips pubip_t
				err = json.Unmarshal(inst.PublicIpAddress, &ips)
				if err == nil {
					pubip = strings.Join(ips.IpAddress, ",")
				}

				state := inst.State.Name
				if state == "" {
					state = fmt.Sprintf("%s", inst.Status)
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
					inst.InstanceId,
					instype,
					inst.InstanceType,
					inst.PublicDnsName,
					pubip,
					inst.PrivateIpAddress,
					strings.ToLower(state))
			}

			d.Info("[instances]")
			w.Flush()

			// separate table for configurations, if any
			if stack.Configuration.Configurations != nil {
				w = tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
				fmt.Fprintf(w, "ROLE\tFLAG\tCONTAINER\n")
				var cnfs []alm.Configurations
				err = json.Unmarshal(stack.Configuration.Configurations, &cnfs)
				cli.ErrorExit(err, 1)

				for _, inst := range cnfs {
					var image string
					if inst.Container != nil {
						contmap := inst.Container.(map[string]interface{})
						for k, v := range contmap {
							if k == "image" {
								image = fmt.Sprintf("%s", v)
							}
						}
					}

					fmt.Fprintf(w, "%s\t%s\t%s\n", inst.Role, inst.Flag, image)
				}

				fmt.Println()
				d.Info("[configurations]")
				w.Flush()
			}
		}
	}
}
