package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
space(s) in it, make sure to surround it with double quotes.`,
	Run: describe,
}

func init() {
	stackCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringP("id", "i", "", "stack id")
	describeCmd.Flags().StringP("fmt", "f", "text", "output format (valid values: text, json)")
	describeCmd.Flags().StringP("out", "o", "", "full file path to write the output")
	describeCmd.Flags().IntP("indent", "n", 2, "indent padding when fmt is 'text' or 'json'")
}

func describe(cmd *cobra.Command, args []string) {
	token, err := util.GetToken()
	if err != nil {
		util.ErrorExit("Cannot read token. See `login` for information on how to login.", 1)
	}

	id := util.GetCliStringFlag(cmd, "id")
	if id == "" {
		util.ErrorExit("Stack id cannot be empty.", 1)
	}

	c := cli.New(util.GetCliStringFlag(cmd, "api-version"))
	ep := c.RootUrl + "/alm/stack/" + fmt.Sprintf("%s", id)
	resp, body, errs := c.GetSafe(ep, fmt.Sprintf("%s", token))
	if errs != nil {
		log.Println("Error(s):", errs)
		os.Exit(1)
	}

	var stacks []stack.DescribeStack
	err = json.Unmarshal(body, &stacks)
	if err != nil {
		log.Println(err)
		var m map[string]interface{}
		err = json.Unmarshal(body, &m)
		if err != nil {
			util.ErrorExit("Internal error.", 1)
		}

		serr := util.ResponseError(resp, m)
		if serr != "" {
			util.ErrorExit(serr, 1)
		}
	}

	switch util.GetCliStringFlag(cmd, "fmt") {
	case "text":
		indent := util.GetCliIntFlag(cmd, "indent")
		stack.PrintR(os.Stdout, &stacks[0], 0, indent)
		f := util.GetCliStringFlag(cmd, "out")
		if f != "" {
			fp, err := os.Create(f)
			if err != nil {
				util.ErrorExit(err.Error(), 1)
			}

			defer fp.Close()
			w := bufio.NewWriter(fp)
			defer w.Flush()
			stack.PrintR(w, &stacks[0], 0, indent)
			log.Println(fmt.Sprintf("Output written to %s.", f))
		}
	case "json":
		indent := util.GetCliIntFlag(cmd, "indent")
		mi, err := json.MarshalIndent(stacks, "", util.Indent(indent))
		if err != nil {
			util.ErrorExit(err.Error(), 1)
		}

		// this should be a prettified JSON output
		fmt.Println(string(mi))

		f := util.GetCliStringFlag(cmd, "out")
		if f != "" {
			err = ioutil.WriteFile(f, mi, 0644)
			if err != nil {
				util.ErrorExit(err.Error(), 1)
			}

			log.Println(fmt.Sprintf("Output written to %s.", f))
		}
	}
}
