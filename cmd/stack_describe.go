package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	term "github.com/buger/goterm"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/stack"
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "display stack details",
	Long: `Display stack details. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

Valid format values: text (default), json, raw, min`,
	Run: describe,
}

func init() {
	stackCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringP("id", "i", "", "stack id")
}

func describe(cmd *cobra.Command, args []string) {
	token, err := util.GetToken()
	if err != nil {
		util.CheckErrorExit("Cannot read token. See `login` for information on how to login.", 1)
	}

	id := util.GetCliStringFlag(cmd, "id")
	if id == "" {
		util.CheckErrorExit("stack id cannot be empty", 1)
	}

	c := cli.New(util.GetCliStringFlag(cmd, "api-version"))
	ep := c.RootUrl + "/alm/stack/" + fmt.Sprintf("%s", id)
	resp, body, errs := c.GetSafe(ep, fmt.Sprintf("%s", token))
	util.CheckErrorExit(errs, 1)

	// we process `--fmt=raw` option first
	out := util.GetCliStringFlag(cmd, "out")
	pfmt := util.GetCliStringFlag(cmd, "fmt")
	if pfmt == "raw" {
		fmt.Println(string(body))
		if out != "" {
			err = util.WriteToFile(out, body)
			util.CheckErrorExit(err, 1)
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
			serr := util.ResponseError(resp, body)
			util.CheckErrorExit(serr, 1)
			util.CheckErrorExit(err, 1)
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
		indent := util.GetCliIntFlag(cmd, "indent")
		mi, err := json.MarshalIndent(sptr, "", util.Indent(indent))
		util.CheckErrorExit(err, 1)
		fmt.Println(string(mi))
		if out != "" {
			err = util.WriteToFile(out, mi)
			util.CheckErrorExit(err, 1)
		}
	default:
		if pfmt == "text" || pfmt == "" {
			indent := util.GetCliIntFlag(cmd, "indent")
			stack.PrintR(os.Stdout, ptr, 0, indent)
			if out != "" {
				fp, err := os.Create(out)
				util.CheckErrorExit(err, 1)

				defer fp.Close()
				w := bufio.NewWriter(fp)
				defer w.Flush()
				stack.PrintR(w, ptr, 0, indent)
				log.Println(fmt.Sprintf("output written to %s", out))
			}
		}
	}
}
