package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/alm/svrconf"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
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
	sid := cli.GetCliStringFlag(cmd, "id")
	if sid == "" {
		cli.ErrorExit("stack id cannot be empty", 1)
	}

	sess, err := clisession()
	cli.ErrorExit(err, 1)

	svc := svrconf.New(sess)
	in := &svrconf.ServerConfigGetInput{
		StackId: cli.GetCliStringFlag(cmd, "id"),
	}

	resp, body, err := svc.Get(in)
	cli.ErrorExit(err, 1)
	exitOn401(resp)

	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")
	if pfmt == "raw" {
		fmt.Println(string(body))
		if out != "" {
			err = ioutil.WriteFile(out, body, 0644)
			cli.ErrorExit(err, 1)
		}

		return
	}

	if pfmt == "json" || pfmt == "" {
		var sc svrconf.ServerConfig
		err = json.Unmarshal(body, &sc)
		cli.ErrorExit(err, 1)

		indent := cli.GetCliIntFlag(cmd, "indent")
		mi, err := json.MarshalIndent(sc, "", pretty.Indent(indent))
		cli.ErrorExit(err, 1)

		fmt.Println(string(mi))

		// write to file option
		out := cli.GetCliStringFlag(cmd, "out")
		if out != "" {
			err = ioutil.WriteFile(out, mi, 0644)
			cli.ErrorExit(err, 1)
		}

		// parse `updated` field for easier reading
		up := time.Unix(sc.Updated, 0)
		d.Info(`parsed value for 'updated' field:`, up.Format(time.RFC1123))
	}
}
