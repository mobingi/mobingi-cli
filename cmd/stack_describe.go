package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

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

Valid format values: text (default), json, raw`,
	Run: describe,
}

func init() {
	stackCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringP("id", "i", "", "stack id")
}

func describe(cmd *cobra.Command, args []string) {
	token, err := util.GetToken()
	if err != nil {
		util.ErrorExit("Cannot read token. See `login` for information on how to login.", 1)
	}

	id := util.GetCliStringFlag(cmd, "id")
	if id == "" {
		util.ErrorExit("stack id cannot be empty", 1)
	}

	c := cli.New(util.GetCliStringFlag(cmd, "api-version"))
	ep := c.RootUrl + "/alm/stack/" + fmt.Sprintf("%s", id)
	resp, body, errs := c.GetSafe(ep, fmt.Sprintf("%s", token))
	if errs != nil {
		log.Println("error(s):", errs)
		os.Exit(1)
	}

	// we process `--fmt=raw` option first
	out := util.GetCliStringFlag(cmd, "out")
	pfmt := util.GetCliStringFlag(cmd, "fmt")
	if pfmt == "raw" {
		fmt.Println(string(body))
		if out != "" {
			err = util.WriteToFile(out, body)
			if err != nil {
				util.ErrorExit(err.Error(), 1)
			}
		}

		return
	}

	// workaround: see description in struct definition
	var ptr interface{}  // pointer to 1st element of slice
	var sptr interface{} // pointer to the whole slice
	var stacks1 []stack.DescribeStack1
	err = json.Unmarshal(body, &stacks1)
	if err != nil {
		var stacks2 []stack.DescribeStack2
		err = json.Unmarshal(body, &stacks2)
		if err != nil {
			serr := util.ResponseError(resp, body)
			if serr != "" {
				util.ErrorExit(serr, 1)
			}

			util.ErrorExit(err.Error(), 1)
		} else {
			ptr = &stacks2[0]
			sptr = stacks2
		}
	} else {
		ptr = &stacks1[0]
		sptr = stacks1
	}

	switch pfmt {
	case "json":
		indent := util.GetCliIntFlag(cmd, "indent")
		mi, err := json.MarshalIndent(sptr, "", util.Indent(indent))
		if err != nil {
			util.ErrorExit(err.Error(), 1)
		}

		fmt.Println(string(mi))
		if out != "" {
			err = util.WriteToFile(out, mi)
			if err != nil {
				util.ErrorExit(err.Error(), 1)
			}
		}
	default:
		if pfmt == "text" || pfmt == "" {
			indent := util.GetCliIntFlag(cmd, "indent")
			stack.PrintR(os.Stdout, ptr, 0, indent)
			if out != "" {
				fp, err := os.Create(out)
				if err != nil {
					util.ErrorExit(err.Error(), 1)
				}

				defer fp.Close()
				w := bufio.NewWriter(fp)
				defer w.Flush()
				stack.PrintR(w, ptr, 0, indent)
				log.Println(fmt.Sprintf("output written to %s", out))
			}
		}
	}
}
