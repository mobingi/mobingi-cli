package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Short: "list all stacks",
	Long: `List all stacks. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

For now, the 'min' format option cannot yet write to a file
using the '--out=[filename]' option. You need to specify either
'text' or 'json'.`,
	Run: list,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("fmt", "f", "min", "output format (valid values: min, text, json)")
	listCmd.Flags().StringP("out", "o", "", "full file path to write the output")
}

func list(cmd *cobra.Command, args []string) {
	token, err := util.GetToken()
	if err != nil {
		util.ErrorExit("Cannot read token. See `login` for information on how to login.", 1)
	}

	c := cli.New(util.GetCliStringFlag(cmd, "api-version"))
	resp, body, errs := c.GetSafe(c.RootUrl+"/alm/stack", fmt.Sprintf("%s", token))
	if errs != nil {
		log.Println("Error(s):", errs)
		os.Exit(1)
	}

	var stacks []stack.ListStack
	err = json.Unmarshal(body, &stacks)
	if err != nil {
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
	case "min":
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

			fmt.Fprintf(stbl, "%s\t%s\t%s\t%s\t%s\t%s\n", s.StackId, s.Nickname, platform, s.StackStatus, s.Configuration.Region, timestr)
		}

		term.Print(stbl)
		term.Flush()
	case "text":
		stack.PrintR(os.Stdout, &stacks[0], 0)
		f := util.GetCliStringFlag(cmd, "out")
		if f != "" {
			fp, err := os.Create(f)
			if err != nil {
				util.ErrorExit(err.Error(), 1)
			}

			defer fp.Close()
			w := bufio.NewWriter(fp)
			defer w.Flush()
			stack.PrintR(w, &stacks[0], 0)
			log.Println(fmt.Sprintf("Output written to %s.", f))
		}
	case "json":
		mi, err := json.MarshalIndent(stacks, "", "  ")
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
