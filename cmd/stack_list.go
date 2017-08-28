package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

func StackListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all stacks",
		Long: `List all stacks. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

Valid format values: min (default), json, raw

For now, the 'min' format option cannot yet write to a file using the '--out=[filename]'
option. You need to specify either 'json', or 'raw'.

Examples:

  $ ` + cmdline.Args0() + ` stack list
  $ ` + cmdline.Args0() + ` stack list --fmt=json --verbose
  $ ` + cmdline.Args0() + ` stack list --fmt=raw --out=/home/foo/tmp.txt`,
		Run: stackList,
	}

	return cmd
}

func stackList(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	d.ErrorExit(err, 1)

	svc := alm.New(sess)
	resp, body, err := svc.List()
	d.ErrorExit(err, 1)
	exitOn401(resp)

	var stacks []alm.ListStack
	err = json.Unmarshal(body, &stacks)
	d.ErrorExit(err, 1)

	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))

		// write to file option
		f := cli.GetCliStringFlag(cmd, "out")
		if f != "" {
			err = ioutil.WriteFile(f, body, 0644)
			d.ErrorExit(err, 1)
			d.Info(fmt.Sprintf("Output written to %s.", f))
		}
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)

		// write to file option
		f := cli.GetCliStringFlag(cmd, "out")
		if f != "" {
			err = ioutil.WriteFile(f, []byte(js), 0644)
			d.ErrorExit(err, 1)
			d.Info(fmt.Sprintf("Output written to %s.", f))
		}
	default:
		if pfmt == "min" || pfmt == "" {
			w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
			fmt.Fprintf(w, "STACK ID\tSTACK NAME\tPLATFORM\tSTATUS\tREGION\tLAUNCHED\n")
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

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					s.StackId,
					s.Nickname,
					platform,
					s.StackStatus,
					s.Configuration.Region,
					timestr)
			}

			w.Flush()
		}
	}
}
