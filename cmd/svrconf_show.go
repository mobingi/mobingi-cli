package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/util"
	"github.com/spf13/cobra"
)

var scshowCmd = &cobra.Command{
	Use:   "show",
	Short: "show current server config",
	Long: `Show current server config. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

Valid format values: json (default), raw`,
	Run: show,
}

func init() {
	svrconfCmd.AddCommand(scshowCmd)
	scshowCmd.Flags().StringP("id", "i", "", "stack id to query")
}

func show(cmd *cobra.Command, args []string) {
	token, err := util.GetToken()
	if err != nil {
		util.ErrorExit("Cannot read token. See `login` for information on how to login.", 1)
	}

	sid := util.GetCliStringFlag(cmd, "id")
	if sid == "" {
		util.ErrorExit("stack id cannot be empty", 1)
	}

	c := cli.New(util.GetCliStringFlag(cmd, "api-version"))
	resp, body, errs := c.GetSafe(c.RootUrl+`/alm/serverconfig?stack_id=`+sid, fmt.Sprintf("%s", token))
	if errs != nil {
		log.Println("Error(s):", errs)
		os.Exit(1)
	}

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

	_ = resp
	if pfmt == "json" || pfmt == "" {
	}

	/*
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
	*/
}
