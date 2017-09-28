package cmd

import (
	"fmt"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/rbac"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

func RbacAttachCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach",
		Short: "attach role to a user",
		Long: `Attach role to a user.

Example:

  $ ` + cmdline.Args0() + ` rbac attach --user foo \
      --role-id morole-58c2297d25645-BtXGMSRsI`,
		Run: rbacAttach,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("user", "", "subuser name")
	cmd.Flags().String("role-id", "", "role id to attach")
	return cmd
}

func rbacAttach(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	cli.ErrorExit(err, 1)

	name := cli.GetCliStringFlag(cmd, "user")
	roleid := cli.GetCliStringFlag(cmd, "role-id")
	in := rbac.AttachRoleToUserInput{
		Username: name,
		RoleId:   roleid,
	}

	svc := rbac.New(sess)
	resp, body, err := svc.AttachRoleToUser(&in)
	cli.ErrorExit(err, 1)
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
