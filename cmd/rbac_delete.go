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

func RbacDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a specific role",
		Long: `Delete a specific role.

Example:

  $ ` + cmdline.Args0() + ` rbac delete --role-id foo`,
		Run: rbacDelete,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().String("role-id", "", "role id to delete")
	return cmd
}

func rbacDelete(cmd *cobra.Command, args []string) {
	sess, err := clisession()
	d.ErrorExit(err, 1)

	svc := rbac.New(sess)

	in := &rbac.DeleteRoleInput{
		RoleId: cli.GetCliStringFlag(cmd, "role-id"),
	}

	resp, body, err := svc.DeleteRole(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	if resp.StatusCode/100 == 2 {
		d.Info(resp.Status)
	} else {
		d.Error(resp.Status)
	}

	fmt.Println(pretty.JSON(string(body), 2))
}
