package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingi/mobingi-cli/pkg/iohelper"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

var showall bool

func TemplateCompareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "compare template versions",
		Long: `Compare template versions.

The options '--src-sid' and '--src-vid' are always required. If '--tgt-sid' is not
provided, it is assumed that target stack id is the same as source stack id. If you
provide '--tgt-vid', '--tgt-body' is ignored. You can get the list of template
versions using the following command ('--id' is the stack id to query):

  $ ` + cmdline.Args0() + ` template versions --id=foo

Examples:

  $ ` + cmdline.Args0() + ` template compare --src-sid=foo \
      --src-vid=bar1 --tgt-vid=bar2,

  $ ` + cmdline.Args0() + ` template compare --src-sid=foo \
      --src-vid=bar --tgt-body=/home/user/tmpl.json`,
		Run: tmplCompare,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("src-sid", "", "source template stack id")
	cmd.Flags().String("src-vid", "", "source template version id")
	cmd.Flags().String("tgt-sid", "", "target template stack id")
	cmd.Flags().String("tgt-vid", "", "target template version id")
	cmd.Flags().String("tgt-body", "", "path to template file to compare to source")
	cmd.Flags().BoolVar(&showall, "show-all", false, "show all results, not just diff")
	return cmd
}

func tmplCompare(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	d.ErrorExit(err, 1)

	svc := alm.New(sess)
	in := &alm.CompareTemplateInput{
		SourceStackId:   cli.GetCliStringFlag(cmd, "src-sid"),
		SourceVersionId: cli.GetCliStringFlag(cmd, "src-vid"),
		TargetStackId:   cli.GetCliStringFlag(cmd, "tgt-sid"),
		TargetVersionId: cli.GetCliStringFlag(cmd, "tgt-vid"),
	}

	tb := cli.GetCliStringFlag(cmd, "tgt-body")
	if tb != "" {
		b, err := ioutil.ReadFile(tb)
		d.ErrorExit(err, 1)

		// set body from file
		in.TargetBody = string(b)
	}

	resp, body, err := svc.CompareTemplate(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	switch pfmt {
	case "raw":
		fmt.Println(string(body))
		if out != "" {
			err = iohelper.WriteToFile(out, body)
			d.ErrorExit(err, 1)
		}
	default:
		if pfmt == "json" || pfmt == "" {
			var tofile string
			if showall {
				indent := cli.GetCliIntFlag(cmd, "indent")
				js := pretty.JSON(string(body), indent)
				fmt.Println(js)
				tofile = js
			} else {
				var m map[string]json.RawMessage
				err = json.Unmarshal(body, &m)
				d.ErrorExit(err, 1)

				diff, ok := m["diff"]
				if ok {
					d.Info("diff:")
					fmt.Println(pretty.JSON(diff, 2))
				} else {
					d.Info("no diff found")
				}
			}

			// write to file option
			if out != "" {
				if tofile != "" {
					err = iohelper.WriteToFile(out, []byte(tofile))
					d.ErrorExit(err, 1)
				} else {
					d.Info("nothing to write to file")
				}
			}
		}
	}
}
