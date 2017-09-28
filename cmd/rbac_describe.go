package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/rbac"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

func RbacDescribeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "describe all defined role(s), or per-user role(s)",
		Long: `List all defined role(s), or per-user role(s).

Examples:

  $ ` + cmdline.Args0() + ` rbac describe
  $ ` + cmdline.Args0() + ` rbac describe --user foo`,
		Run: rbacDescribe,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("user", "", "subuser name, all if empty")
	return cmd
}

func rbacDescribe(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	cli.ErrorExit(err, 1)

	svc := rbac.New(sess)

	var in *rbac.DescribeRolesInput
	user := cli.GetCliStringFlag(cmd, "user")
	if user != "" {
		in = &rbac.DescribeRolesInput{User: user}
	}

	resp, body, err := svc.DescribeRoles(in)
	cli.ErrorExit(err, 1)
	exitOn401(resp)

	if resp.StatusCode/100 != 2 {
		d.Error(resp.Status)
	} else {
		d.Info(resp.Status)
	}

	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")

	var outb []byte

	switch pfmt {
	case "raw":
		fmt.Println(string(body))
		outb = body
	case "json":
		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)
		outb = []byte(js)
	default:
		if pfmt == "min" || pfmt == "" {
			if resp.StatusCode/100 != 2 {
				// error already printed
				os.Exit(1)
			}

			var rj []json.RawMessage
			err = json.Unmarshal(body, &rj)
			cli.ErrorExit(err, 1)

			w := tabwriter.NewWriter(os.Stdout, 0, 10, 5, ' ', 0)
			fmt.Fprintf(w, "ROLE ID\tUSER ID\tNAME\tCREATE TIME\tUPDATE TIME\tSCOPE\n")
			var roleid, userid, name, ct, ut, scope string
			for _, item := range rj {
				// var m map[string]interface{}
				var m map[string]json.RawMessage
				err = json.Unmarshal(item, &m)
				cli.ErrorExit(err, 1)

				_roleid, ok := m["role_id"]
				if ok {
					roleid = fmt.Sprintf("%s", _roleid)
					roleid = strings.Trim(roleid, "\"")
				}

				_userid, ok := m["user_id"]
				if ok {
					userid = fmt.Sprintf("%s", _userid)
					userid = strings.Trim(userid, "\"")
				}

				_name, ok := m["name"]
				if ok {
					name = fmt.Sprintf("%s", _name)
					name = strings.Trim(name, "\"")
				}

				ct = "-"
				_ct, ok := m["create_time"]
				if ok {
					_cts := strings.Trim(string(_ct), "\"")
					t, err := time.Parse(time.RFC3339, _cts)
					if err == nil {
						ct = t.Format(time.RFC1123)
					}
				}

				ut = "-"
				_ut, ok := m["create_time"]
				if ok {
					_uts := strings.Trim(string(_ut), "\"")
					t, err := time.Parse(time.RFC3339, _uts)
					if err == nil {
						ut = t.Format(time.RFC1123)
					}
				}

				_scope, ok := m["scope"]
				if ok {
					scope = string(_scope)
					scope = strings.Trim(scope, "\"")
					if len(scope) >= 23 {
						scope = scope[:23] + "..."
					}
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					roleid,
					userid,
					name,
					ct,
					ut,
					scope)
			}

			w.Flush()
		}
	}

	if out != "" {
		err = ioutil.WriteFile(out, outb, 0644)
		cli.ErrorExit(err, 1)
	}
}
