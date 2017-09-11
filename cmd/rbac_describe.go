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
	d.ErrorExit(err, 1)

	svc := rbac.New(sess)

	var in *rbac.DescribeRolesInput
	user := cli.GetCliStringFlag(cmd, "user")
	if user != "" {
		in = &rbac.DescribeRolesInput{User: user}
	}

	resp, body, err := svc.DescribeRoles(in)
	d.ErrorExit(err, 1)
	exitOn401(resp)

	out := cli.GetCliStringFlag(cmd, "out")
	pfmt := cli.GetCliStringFlag(cmd, "fmt")

	var outb []byte

	switch pfmt {
	case "raw":
		fmt.Println(string(body))
		outb = body
	default:
		indent := cli.GetCliIntFlag(cmd, "indent")
		js := pretty.JSON(string(body), indent)
		fmt.Println(js)
		outb = []byte(js)
	}

	if out != "" {
		err = ioutil.WriteFile(out, outb, 0644)
		d.ErrorExit(err, 1)
	}
}
