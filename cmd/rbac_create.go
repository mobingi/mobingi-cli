package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/rbac"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

var roleAllowAll bool

func RbacCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "define a role or create a user role",
		Long: `Define a role or create a user role.

Example:

  $ ` + cmdline.Args0() + ` rbac create --name testrole \
      --scope /home/user/role.json`,
		Run: rbacCreate,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("type", "role", "create type: role, user")
	cmd.Flags().String("name", "", "role name (when type is role)")
	cmd.Flags().String("scope", "", "path to role file (when type is role)")
	cmd.Flags().BoolVar(&roleAllowAll, "allow-all", false, "true if you allow all access")
	return cmd
}

func rbacCreate(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	d.ErrorExit(err, 1)

	var in rbac.CreateRoleInput

	typ := cli.GetCliStringFlag(cmd, "type")
	name := cli.GetCliStringFlag(cmd, "name")
	scope := cli.GetCliStringFlag(cmd, "scope")

	switch typ {
	case "role":
		if roleAllowAll {
			in = rbac.CreateRoleInput{
				Name:  name,
				Scope: *(rbac.NewRoleAll("Allow")),
			}
		} else {
			b, err := ioutil.ReadFile(scope)
			d.ErrorExit(err, 1)

			var rr rbac.Role
			err = json.Unmarshal(b, &rr)
			d.ErrorExit(err, 1)

			in = rbac.CreateRoleInput{
				Name:  name,
				Scope: rr,
			}
		}
	case "user":
		d.Info("not yet supported")
		os.Exit(0)
	default:
		d.ErrorExit("unknown type", 1)
	}

	svc := rbac.New(sess)
	resp, body, err := svc.CreateRole(&in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	indent := cli.GetCliIntFlag(cmd, "indent")
	js := pretty.JSON(string(body), indent)
	if resp.StatusCode/100 != 2 {
		d.Error(resp.Status)
	} else {
		d.Info(resp.Status)
	}

	fmt.Println(js)
}
