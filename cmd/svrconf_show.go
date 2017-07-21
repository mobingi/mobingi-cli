package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mobingilabs/mocli/api"
	"github.com/mobingilabs/mocli/pkg/check"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/svrconf"
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
	var err error
	sid := util.GetCliStringFlag(cmd, "id")
	if sid == "" {
		check.ErrorExit("stack id cannot be empty", 1)
	}

	c := api.NewClient(api.NewConfig(cmd))
	resp, body, errs := c.Get(`/alm/serverconfig?stack_id=` + sid)
	check.ErrorExit(errs, 1)

	out := util.GetCliStringFlag(cmd, "out")
	pfmt := util.GetCliStringFlag(cmd, "fmt")
	if pfmt == "raw" {
		fmt.Println(string(body))
		if out != "" {
			err = util.WriteToFile(out, body)
			check.ErrorExit(err, 1)
		}

		return
	}

	if pfmt == "json" || pfmt == "" {
		var sc svrconf.ServerConfig
		err = json.Unmarshal(body, &sc)
		check.ErrorExit(err, 1)

		serr := check.ResponseError(resp, body)
		check.ErrorExit(serr, 1)

		indent := util.GetCliIntFlag(cmd, "indent")
		mi, err := json.MarshalIndent(sc, "", util.Indent(indent))
		check.ErrorExit(err, 1)

		fmt.Println(string(mi))

		// write to file option
		out := util.GetCliStringFlag(cmd, "out")
		if out != "" {
			err = util.WriteToFile(out, mi)
			check.ErrorExit(err, 1)
		}

		// parse `updated` field for easier reading
		up := time.Unix(sc.Updated, 0)
		d.Info(`parsed value for 'updated' field:`, up.Format(time.RFC1123))
	}
}
