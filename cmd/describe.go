package cmd

import (
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
	Long:  `Display stack details.`,
	Run:   describe,
}

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.Flags().String("id", "", "stack id")
}

func describe(cmd *cobra.Command, args []string) {
	token, err := util.GetToken()
	if err != nil {
		util.PrintErrorAndExit("Cannot read token. See `login` for information on how to login.", 1)
	}

	id := util.GetCliStringFlag(cmd, "id")
	if id == "" {
		util.PrintErrorAndExit("Stack id cannot be empty.", 1)
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
			util.PrintErrorAndExit("Internal error.", 1)
		}

		serr := util.BuildRequestError(resp, m)
		if serr != "" {
			util.PrintErrorAndExit(serr, 1)
		}
	}

	stbl := term.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintf(stbl, "INSTANCE ID\tINSTANCE TYPE\tPUBLIC IP\tPRIVATE IP\n")
	for _, s := range stacks {
		for _, i := range s.Instances {
			fmt.Fprintf(stbl, "%s\t%s\t%s\t%s\n", i.InstanceId, i.InstanceType, i.PublicIpAddress, i.PrivateIpAddress)
		}
	}

	term.Print(stbl)
	term.Flush()
}
