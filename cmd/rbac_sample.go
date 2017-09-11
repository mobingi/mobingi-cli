package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/mobingi/mobingi-cli/pkg/cli"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/rbac"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/cmdline"
	d "github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/pretty"
	"github.com/spf13/cobra"
)

func RbacSampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sample",
		Short: "print a sample role",
		Long: `Print a sample role.

Example:

  $ ` + cmdline.Args0() + ` rbac sample
  $ ` + cmdline.Args0() + ` rbac sample --out=/home/user/sample.json`,
		Run: rbacSample,
	}

	return cmd
}

func rbacSample(cmd *cobra.Command, args []string) {
	sample := rbac.Role{}
	sample.Version = "2017-05-05"
	sample.Statement = make([]rbac.RoleStatement, 0)
	sample.Statement = append(sample.Statement, rbac.RoleStatement{
		Effect:   "Deny",
		Action:   []string{"stack:describeStacks"},
		Resource: []string{"mrn:alm:stack:mo-xxxxxxx"},
	})

	sample.Statement = append(sample.Statement, (*(rbac.NewRoleAll("Allow"))).Statement[0])
	str := pretty.JSON(sample, 2)
	fmt.Println(str)

	out := cli.GetCliStringFlag(cmd, "out")
	if out != "" {
		err := ioutil.WriteFile(out, []byte(str), 0644)
		d.ErrorExit(err, 1)
		d.Info("sample written to", out)
	}
}
