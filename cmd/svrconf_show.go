package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mobingilabs/mocli/pkg/cli"
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

	if pfmt == "json" || pfmt == "" {
		var sc svrconf.ServerConfig
		err = json.Unmarshal(body, &sc)
		if err != nil {
			log.Println(err)
			var m map[string]interface{}
			err = json.Unmarshal(body, &m)
			if err != nil {
				util.ErrorExit("internal error", 1)
			}

			serr := util.ResponseError(resp, m)
			if serr != "" {
				util.ErrorExit(serr, 1)
			}
		}

		indent := util.GetCliIntFlag(cmd, "indent")
		mi, err := json.MarshalIndent(sc, "", util.Indent(indent))
		if err != nil {
			util.ErrorExit(err.Error(), 1)
		}

		fmt.Println(string(mi))
		out := util.GetCliStringFlag(cmd, "out")
		if out != "" {
			err = util.WriteToFile(out, mi)
			if err != nil {
				util.ErrorExit(err.Error(), 1)
			}
		}

		// extra information on `updated` field
		up := time.Unix(sc.Updated, 0)
		log.Println("updated (parsed):", up.Format(time.RFC1123))
	}
}
