package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	term "github.com/buger/goterm"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/stack"
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all stack",
	Long:  `List all stack.`,
	Run:   list,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) {
	token, err := util.GetToken()
	if err != nil {
		util.PrintErrorAndExit("Cannot read token. See `login` for information on how to login.", 1)
	}

	c := cli.New(util.GetCliStringFlag(cmd, "api-version"))
	resp, body, errs := c.GetSafe(c.RootUrl+"/alm/stack", fmt.Sprintf("%s", token))
	if errs != nil {
		log.Println("Error(s):", errs)
		os.Exit(1)
	}

	var stacks []stack.Stack
	err = json.Unmarshal(body, &stacks)
	if err != nil {
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
	fmt.Fprintf(stbl, "STACK ID\tSTACK NAME\tPLATFORM\tSTATUS\tREGION\tLAUNCHED\n")
	for _, s := range stacks {
		timestr := s.CreateTime
		t, err := time.Parse(time.RFC3339, s.CreateTime)
		if err == nil {
			timestr = t.Format(time.RFC1123)
		}

		platform := "?"
		if s.Configuration.Aws != "" {
			platform = "AWS"
		}

		fmt.Fprintf(stbl, "%s\t%s\t%s\t%s\t%s\t%s\n", s.StackId, s.Nickname, platform, s.StackStatus, s.Configuration.Region, timestr)
	}

	term.Print(stbl)
	term.Flush()
	// log.Println(string(body))
}
