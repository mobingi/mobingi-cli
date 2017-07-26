package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/mobingilabs/mocli/pkg/iohelper"
	"github.com/mobingilabs/mocli/pkg/pretty"
	"github.com/mobingilabs/mocli/pkg/svrconf"
	"github.com/spf13/cobra"
)

func ServerConfigShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show current server config",
		Long: `Show current server config. If you specify the '--out=[filename]' option,
make sure you provide the full path of the file. If the path has
space(s) in it, make sure to surround it with double quotes.

Valid format values: json (default), raw`,
		Run: show,
	}

	cmd.Flags().StringP("id", "i", "", "stack id to query")
	return cmd
}

func show(cmd *cobra.Command, args []string) {
	var err error
	sid := cli.GetCliStringFlag(cmd, "id")
	if sid == "" {
		check.ErrorExit("stack id cannot be empty", 1)
	}

	c := client.NewClient(client.NewApiConfig(cmd))
	body, err := c.AuthGet(`/alm/serverconfig?stack_id=` + sid)
	check.ErrorExit(err, 1)

	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	if pfmt == "raw" {
		fmt.Println(string(body))
		if out != "" {
			err = iohelper.WriteToFile(out, body)
			check.ErrorExit(err, 1)
		}

		return
	}

	if pfmt == "json" || pfmt == "" {
		var sc svrconf.ServerConfig
		err = json.Unmarshal(body, &sc)
		check.ErrorExit(err, 1)

		indent := cli.GetCliIntFlag(cmd, "indent")
		mi, err := json.MarshalIndent(sc, "", pretty.Indent(indent))
		check.ErrorExit(err, 1)

		fmt.Println(string(mi))

		// write to file option
		out := cli.GetCliStringFlag(cmd, "out")
		if out != "" {
			err = iohelper.WriteToFile(out, mi)
			check.ErrorExit(err, 1)
		}

		// parse `updated` field for easier reading
		up := time.Unix(sc.Updated, 0)
		d.Info(`parsed value for 'updated' field:`, up.Format(time.RFC1123))
	}
}
